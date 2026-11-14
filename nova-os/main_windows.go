//go:build windows
package main

import (
  "bytes"
  _ "embed"
  "errors"
  "fmt"
  "io"
  "net/url"
  "os"
  "os/exec"
  "path/filepath"
  "strings"
  "syscall"
  "unicode/utf16"
  "unsafe"
)

//go:embed index.html
var uiHTML []byte
//go:embed icon.ico
var appIcon []byte

const appName = "NOVA OS"

func wide(s string) *uint16 {
 u := utf16.Encode([]rune(s+"\x00"))
 return &u[0]
}
var msgbox=syscall.NewLazyDLL("user32.dll").NewProc("MessageBoxW")
func notify(text string, kind uintptr) uintptr {
 ret,_,_:=msgbox.Call(0,uintptr(unsafe.Pointer(wide(text))),uintptr(unsafe.Pointer(wide(appName))),kind)
 return ret
}
func appData() (string,error){
 root:=os.Getenv("LOCALAPPDATA")
 if root=="" {return "",errors.New("LOCALAPPDATA dizini bulunamadı")}
 return filepath.Join(root,"NOVA OS"),nil
}
func writeIfDifferent(path string, b []byte) error {
 old,e:=os.ReadFile(path); if e==nil && bytes.Equal(old,b) {return nil}
 if err:=os.MkdirAll(filepath.Dir(path),0700);err!=nil{return err}
 tmp:=path+".new"
 if err:=os.WriteFile(tmp,b,0600);err!=nil{return err}
 if err:=os.Rename(tmp,path);err!=nil {os.Remove(path);return os.Rename(tmp,path)}
 return nil
}
func copyFile(from,to string)error{
 f,err:=os.Open(from);if err!=nil{return err};defer f.Close()
 if err:=os.MkdirAll(filepath.Dir(to),0700);err!=nil{return err}
 out,err:=os.Create(to);if err!=nil{return err}
 _,err=io.Copy(out,f); if e:=out.Close();err==nil{err=e}; return err
}
func findEdge()string{
 var list []string
 for _,key:=range []string{"ProgramFiles(x86)","ProgramFiles","LOCALAPPDATA"}{
  if p:=os.Getenv(key);p!=""{list=append(list,filepath.Join(p,"Microsoft","Edge","Application","msedge.exe"))}
 }
 for _,p:=range list{if fi,e:=os.Stat(p);e==nil&&!fi.IsDir(){return p}}
 if p,e:=exec.LookPath("msedge.exe");e==nil{return p}
 return ""
}
func fileURL(path string)string{
 abs,_:=filepath.Abs(path)
 p:=filepath.ToSlash(abs)
 if !strings.HasPrefix(p,"/"){p="/"+p}
 return (&url.URL{Scheme:"file",Path:p}).String()
}
func launch()error{
 root,err:=appData();if err!=nil{return err}
 path:=filepath.Join(root,"UI","index.html")
 if err:=writeIfDifferent(path,uiHTML);err!=nil{return fmt.Errorf("Arayüz oluşturulamadı: %w",err)}
 iconPath:=filepath.Join(root,"UI","nova.ico")
 _=writeIfDifferent(iconPath,appIcon)
 edge:=findEdge()
 if edge==""{return errors.New("Microsoft Edge bulunamadı. NOVA OS, Windows 11'de kurulu olan Edge'in uygulama motorunu kullanır. Edge kurulu olmalıdır.")}
 args:=[]string{"--app="+fileURL(path),"--user-data-dir="+filepath.Join(root,"Profile"),"--no-first-run","--disable-extensions","--window-size=1440,900"}
 cmd:=exec.Command(edge,args...)
 cmd.SysProcAttr=&syscall.SysProcAttr{HideWindow:true}
 if err:=cmd.Start();err!=nil{return fmt.Errorf("NOVA OS başlatılamadı: %w",err)}
 return nil
}
func vbquote(s string)string{return strings.ReplaceAll(s,"\"","\"\"")}
func shortcut(dest,target,icon string)error{
 body:=fmt.Sprintf("Set sh = CreateObject(\"WScript.Shell\")\r\nSet lnk = sh.CreateShortcut(\"%s\")\r\nlnk.TargetPath = \"%s\"\r\nlnk.WorkingDirectory = \"%s\"\r\nlnk.IconLocation = \"%s,0\"\r\nlnk.Description = \"NOVA OS - Hafif sanal masaustu\"\r\nlnk.Save\r\n",vbquote(dest),vbquote(target),vbquote(filepath.Dir(target)),vbquote(icon))
 f,err:=os.CreateTemp("","nova-shortcut-*.vbs");if err!=nil{return err}
 tmp:=f.Name();defer os.Remove(tmp)
 if _,err:=f.WriteString(body);err!=nil{f.Close();return err};f.Close()
 c:=exec.Command("wscript.exe","//B","//Nologo",tmp)
 c.SysProcAttr=&syscall.SysProcAttr{HideWindow:true}
 return c.Run()
}
func install() error {
 root,err:=appData();if err!=nil{return err}
 program:=filepath.Join(os.Getenv("LOCALAPPDATA"),"Programs","NOVA OS")
 dest:=filepath.Join(program,"NOVA OS.exe")
 src,err:=os.Executable();if err!=nil{return err}
 if !strings.EqualFold(filepath.Clean(src),filepath.Clean(dest)){
  if err:=copyFile(src,dest);err!=nil{return fmt.Errorf("EXE kurulamadı: %w",err)}
 }
 icon:=filepath.Join(root,"UI","nova.ico")
 if err:=writeIfDifferent(icon,appIcon);err!=nil{return err}
 start:=filepath.Join(os.Getenv("APPDATA"),"Microsoft","Windows","Start Menu","Programs")
 desktop:=filepath.Join(os.Getenv("USERPROFILE"),"Desktop")
 if err:=os.MkdirAll(start,0700);err==nil { _=shortcut(filepath.Join(start,"NOVA OS.lnk"),dest,icon) }
 if _,e:=os.Stat(desktop);e==nil { _=shortcut(filepath.Join(desktop,"NOVA OS.lnk"),dest,icon) }
 notify("NOVA OS kuruldu!\n\nBaşlat menüsüne ve masaüstüne kısayol eklendi.\nDosyaların ve ayarların NOVA OS profilinde saklanır.\n\nSistem şimdi açılacak.",0x40)
 return launch()
}
func main(){
 exe,_:=os.Executable()
 setup:=strings.Contains(strings.ToLower(filepath.Base(exe)),"setup")
 for _,a:=range os.Args[1:] { if a=="--install"{setup=true};if a=="--portable"{setup=false} }
 var err error
 if setup {err=install()} else {err=launch()}
 if err!=nil{notify("İşlem tamamlanamadı:\n\n"+err.Error(),0x10)}
}