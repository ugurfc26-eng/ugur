//go:build windows
// NOVA OS Web Installer v0.2 - downloads only a named release asset from GitHub.
package main

import (
 "archive/zip"
 "context"
 "crypto/sha256"
 "encoding/hex"
 _ "embed"
 "errors"
 "fmt"
 "io"
 "net/http"
 "net/url"
 "os"
 "os/exec"
 "path/filepath"
 "runtime"
 "strconv"
 "strings"
 "syscall"
 "time"
 "unsafe"
)

const (
 releasePage="https://github.com/ugurfc26-eng/ugur/releases/tag/nova-os-v0.2.0"
 releaseBase="https://github.com/ugurfc26-eng/ugur/releases/download/nova-os-v0.2.0/"
 payloadName="NOVA_OS_Package.zip"
 hashName="SHA256SUMS.txt"
 WM_DESTROY=0x0002
 WM_COMMAND=0x0111
 WM_CLOSE=0x0010
 WS_VISIBLE=0x10000000
 WS_CHILD=0x40000000
 WS_OVERLAPPED=0x00000000
 WS_CAPTION=0x00C00000
 WS_SYSMENU=0x00080000
 WS_MINIMIZEBOX=0x00020000
 WS_TABSTOP=0x00010000
 WS_EX_APPWINDOW=0x00040000
 BN_CLICKED=0
 ID_INSTALL=1001
 ID_SITE=1002
 ID_CANCEL=1003
 SW_SHOWNORMAL=1
 MB_OK=0x0
 MB_ICONERROR=0x10
 MB_ICONINFORMATION=0x40
)

//go:embed fallback.zip
var fallbackPackage []byte

var (
 user32=syscall.NewLazyDLL("user32.dll")
 kernel32=syscall.NewLazyDLL("kernel32.dll")
 shell32=syscall.NewLazyDLL("shell32.dll")
 procRegisterClassEx=user32.NewProc("RegisterClassExW")
 procCreateWindow=user32.NewProc("CreateWindowExW")
 procDefWindow=user32.NewProc("DefWindowProcW")
 procShowWindow=user32.NewProc("ShowWindow")
 procUpdateWindow=user32.NewProc("UpdateWindow")
 procGetMessage=user32.NewProc("GetMessageW")
 procTranslateMessage=user32.NewProc("TranslateMessage")
 procDispatchMessage=user32.NewProc("DispatchMessageW")
 procDestroyWindow=user32.NewProc("DestroyWindow")
 procPostQuit=user32.NewProc("PostQuitMessage")
 procSetWindowText=user32.NewProc("SetWindowTextW")
 procEnableWindow=user32.NewProc("EnableWindow")
 procMessageBox=user32.NewProc("MessageBoxW")
 procShellExecute=shell32.NewProc("ShellExecuteW")
 procGetModuleHandle=kernel32.NewProc("GetModuleHandleW")
 mainWindow uintptr
 installButton uintptr
 siteButton uintptr
 statusLabel uintptr
 progressLabel uintptr
 busy bool
)

type WndClassEx struct {
 Size uint32
 Style uint32
 WndProc uintptr
 ClsExtra int32
 WndExtra int32
 Instance uintptr
 Icon uintptr
 Cursor uintptr
 Background uintptr
 MenuName *uint16
 ClassName *uint16
 IconSmall uintptr
}
type POINT struct{X,Y int32}
type MSG struct {Hwnd uintptr;Message uint32;_ uint32;WParam uintptr;LParam uintptr;Time uint32;Pt POINT;LPrivate uint32;_ uint32}
func w(s string)*uint16{p,_:=syscall.UTF16PtrFromString(s);return p}
func setText(handle uintptr,s string){procSetWindowText.Call(handle,uintptr(unsafe.Pointer(w(s))))}
func notify(s string,kind uintptr){procMessageBox.Call(mainWindow,uintptr(unsafe.Pointer(w(s))),uintptr(unsafe.Pointer(w("NOVA OS Kurulum"))),kind)}
func showSite(){procShellExecute.Call(mainWindow,uintptr(unsafe.Pointer(w("open"))),uintptr(unsafe.Pointer(w(releasePage))),0,0,SW_SHOWNORMAL)}
func child(class,name string,style uint32,x,y,width,height int32,id uintptr)uintptr{
 h,_,_:=procCreateWindow.Call(0,uintptr(unsafe.Pointer(w(class))),uintptr(unsafe.Pointer(w(name))),uintptr(style|WS_CHILD|WS_VISIBLE),uintptr(x),uintptr(y),uintptr(width),uintptr(height),mainWindow,id,0,0)
 return h
}
func wndProc(hwnd uintptr,msg uint32,wp,lp uintptr) uintptr {
 switch msg {
 case WM_COMMAND:
  id:=wp&0xffff
  switch id {
  case ID_INSTALL:
   if !busy {busy=true;procEnableWindow.Call(installButton,0);go doInstall()}
  case ID_SITE:showSite()
  case ID_CANCEL:if !busy {procDestroyWindow.Call(hwnd)}
  }
  return 0
 case WM_CLOSE:
  if !busy {procDestroyWindow.Call(hwnd)}
  return 0
 case WM_DESTROY:
  procPostQuit.Call(0);return 0
 }
 r,_,_:=procDefWindow.Call(hwnd,uintptr(msg),wp,lp);return r
}
func createWindow()error{
 inst,_,_:=procGetModuleHandle.Call(0)
 cls:=w("NovaOSWebInstaller02")
 wnd:=WndClassEx{Size:uint32(unsafe.Sizeof(WndClassEx{})),WndProc:syscall.NewCallback(wndProc),Instance:inst,Background:6,ClassName:cls}
 if n,_,_:=procRegisterClassEx.Call(uintptr(unsafe.Pointer(&wnd)));n==0{return errors.New("Pencere sınıfı oluşturulamadı")}
 mainWindow,_,_=procCreateWindow.Call(WS_EX_APPWINDOW,uintptr(unsafe.Pointer(cls)),uintptr(unsafe.Pointer(w("NOVA OS • Kurulum Merkezi"))),WS_OVERLAPPED|WS_CAPTION|WS_SYSMENU|WS_MINIMIZEBOX,300,180,580,375,0,0,inst,0)
 if mainWindow==0{return errors.New("Kurulum penceresi oluşturulamadı")}
 child("STATIC","NOVA OS  |  WEB KURULUM",0,26,21,500,34,0)
 child("STATIC","Hafif sanal masaüstü • Windows 11 x64 • VirtualBox gerektirmez",0,26,58,520,25,0)
 child("STATIC","KUR VE İNDİR: Güvenli GitHub yayınından paketi indirir ve kurar.",0,26,100,525,35,0)
 child("STATIC","Kurulum kişisel dosyalarını silmez. İnternet bağlantısı gerekir.",0,26,135,520,30,0)
 statusLabel=child("STATIC","Hazır. Aşağıdan KUR VE İNDİR düğmesine bas.",0,26,178,530,26,0)
 progressLabel=child("STATIC","",0,26,209,530,26,0)
 installButton=child("BUTTON","KUR VE İNDİR",WS_TABSTOP,26,260,204,42,ID_INSTALL)
 siteButton=child("BUTTON","İndirme sayfasını aç",WS_TABSTOP,241,260,185,42,ID_SITE)
 child("BUTTON","Kapat",WS_TABSTOP,436,260,107,42,ID_CANCEL)
 procShowWindow.Call(mainWindow,1);procUpdateWindow.Call(mainWindow)
 var msg MSG
 for {v,_,_:=procGetMessage.Call(uintptr(unsafe.Pointer(&msg)),0,0,0);if int32(v)<=0 {break};procTranslateMessage.Call(uintptr(unsafe.Pointer(&msg)));procDispatchMessage.Call(uintptr(unsafe.Pointer(&msg)))}
 return nil
}
func downloadable() *http.Client{
 return &http.Client{Timeout:0,CheckRedirect:func(req *http.Request,via []*http.Request)error{
  if req.URL.Scheme!="https" {return errors.New("Güvensiz indirme bağlantısı reddedildi")}
  h:=strings.ToLower(req.URL.Hostname())
  allowed:=h=="github.com"||h=="objects.githubusercontent.com"||h=="release-assets.githubusercontent.com"||h=="api.github.com"
  if !allowed {return errors.New("Bilinmeyen indirme sunucusu: "+h)}
  if len(via)>12{return errors.New("Çok fazla yönlendirme")};return nil
 }}
}
func download(ctx context.Context,c *http.Client,remote,local string,max int64,progress bool)error{
 u,e:=url.Parse(remote);if e!=nil||u.Scheme!="https"||u.Hostname()!="github.com"{return errors.New("İndirme adresi doğrulanamadı")}
 req,e:=http.NewRequestWithContext(ctx,"GET",remote,nil);if e!=nil{return e}
 req.Header.Set("User-Agent","NOVA-OS-Installer/0.2")
 resp,e:=c.Do(req);if e!=nil{return e};defer resp.Body.Close()
 if resp.StatusCode!=http.StatusOK {return fmt.Errorf("GitHub HTTP %d; yayın henüz hazır olmayabilir",resp.StatusCode)}
 if resp.ContentLength>max {return errors.New("Paket beklenenden büyük")}
 out,e:=os.Create(local);if e!=nil{return e};defer out.Close()
 reader:=io.LimitReader(resp.Body,max+1)
 buf:=make([]byte,64*1024);var n int64;last:=-1
 for {k,er:=reader.Read(buf);if k>0 {if _,werr:=out.Write(buf[:k]);werr!=nil{return werr};n+=int64(k);if n>max{return errors.New("Maksimum indirme boyutu aşıldı")};if progress&&resp.ContentLength>0 {pct:=int(n*100/resp.ContentLength);if pct!=last&&(pct%5==0||pct==100){last=pct;setText(progressLabel,"İndiriliyor: %"+strconv.Itoa(pct)+" • "+fmt.Sprintf("%.1f MB",float64(n)/(1024*1024)))}}};if er==io.EOF{break};if er!=nil{return er}}
 return out.Sync()
}
func expectedSHA(s,filename string)(string,error){
 for _,line:=range strings.Split(s,"\n"){
  fields:=strings.Fields(line)
  if len(fields)>=2&&strings.TrimPrefix(fields[1],"*")==filename&&len(fields[0])==64 {if _,e:=hex.DecodeString(fields[0]);e==nil{return strings.ToLower(fields[0]),nil}}
 }
 return "",errors.New("SHA256 imzası bulunamadı")
}
func verifyHash(file,expected string)error{
 f,e:=os.Open(file);if e!=nil{return e};defer f.Close();h:=sha256.New();if _,e=io.Copy(h,f);e!=nil{return e};if hex.EncodeToString(h.Sum(nil))!=expected{return errors.New("İndirilen paket SHA256 kontrolünden geçmedi")};return nil
}
func extractExecutable(zippath,target string)error{
 f,e:=os.Open(zippath);if e!=nil{return e};defer f.Close()
 fi,e:=f.Stat();if e!=nil{return e}
 zr,e:=zip.NewReader(f,fi.Size());if e!=nil{return e}
 if len(zr.File)>15{return errors.New("Paket yapısı beklenmiyor")}
 for _,entry:=range zr.File {
  if entry.Name!="NOVA_OS.exe" {continue}
  if entry.UncompressedSize64>200<<20{return errors.New("EXE boyutu çok büyük")}
  r,e:=entry.Open();if e!=nil{return e};defer r.Close()
  if e=os.MkdirAll(filepath.Dir(target),0700);e!=nil{return e}
  tmp:=target+".new";out,e:=os.OpenFile(tmp,os.O_CREATE|os.O_TRUNC|os.O_WRONLY,0700);if e!=nil{return e}
  n,e:=io.Copy(out,io.LimitReader(r,200<<20+1));closeE:=out.Close();if e!=nil{os.Remove(tmp);return e};if closeE!=nil{os.Remove(tmp);return closeE};if n!=int64(entry.UncompressedSize64){os.Remove(tmp);return errors.New("Paket eksik açıldı")}
  if e=os.Rename(tmp,target);e!=nil { // Windows cannot always rename over an older installed EXE.
   backup:=target+".old";_ =os.Remove(backup)
   if old:=os.Rename(target,backup);old!=nil{os.Remove(tmp);return e}
   if e=os.Rename(tmp,target);e!=nil{_ = os.Rename(backup,target);return e}
   _ =os.Remove(backup)
  }
  return nil
 }
 return errors.New("NOVA_OS.exe pakette yok")
}
func shortcut(dst,target string)error{
 // Windows shortcuts are created only after the user requests installation.
 body:=fmt.Sprintf("Set sh = CreateObject(\"WScript.Shell\")\r\nSet lnk = sh.CreateShortcut(\"%s\")\r\nlnk.TargetPath = \"%s\"\r\nlnk.WorkingDirectory = \"%s\"\r\nlnk.IconLocation = \"%s,0\"\r\nlnk.Description = \"NOVA OS sanal masaüstü\"\r\nlnk.Save\r\n",strings.ReplaceAll(dst,"\"","\"\""),strings.ReplaceAll(target,"\"","\"\""),strings.ReplaceAll(filepath.Dir(target),"\"","\"\""),strings.ReplaceAll(target,"\"","\"\""))
 file,e:=os.CreateTemp("","nova-shortcut-*.vbs");if e!=nil{return e};defer os.Remove(file.Name());if _,e=file.WriteString(body);e!=nil{file.Close();return e};file.Close();cmd:=exec.Command("wscript.exe","//B","//Nologo",file.Name());cmd.SysProcAttr=&syscall.SysProcAttr{HideWindow:true};return cmd.Run()
}
func install(ctx context.Context)error{
 root:=os.Getenv("LOCALAPPDATA");if root==""{return errors.New("LOCALAPPDATA bulunamadı")}
 temp,e:=os.MkdirTemp("","nova-webinstall-*");if e!=nil{return e};defer os.RemoveAll(temp)
 c:=downloadable()
 hashfile:=filepath.Join(temp,"sha.txt")
 setText(statusLabel,"Paket doğrulama dosyası indiriliyor...")
 if e=download(ctx,c,releaseBase+hashName,hashfile,16*1024,false);e!=nil{return e}
 sumData,e:=os.ReadFile(hashfile);if e!=nil{return e};hash,e:=expectedSHA(string(sumData),payloadName);if e!=nil{return e}
 payload:=filepath.Join(temp,"package.zip")
 setText(statusLabel,"NOVA OS paketi GitHub'dan indiriliyor...")
 if e=download(ctx,c,releaseBase+payloadName,payload,500<<20,true);e!=nil{return e}
 setText(statusLabel,"SHA256 dosya bütünlüğü denetleniyor...")
 if e=verifyHash(payload,hash);e!=nil{return e}
 return installZip(payload)
}
func installZip(payload string)error{
 root:=os.Getenv("LOCALAPPDATA");if root==""{return errors.New("LOCALAPPDATA bulunamadı")}
 dest:=filepath.Join(root,"Programs","NOVA OS","NOVA OS.exe")
 setText(statusLabel,"Uygulama kuruluyor...")
 if e:=extractExecutable(payload,dest);e!=nil{return e}
 start:=filepath.Join(os.Getenv("APPDATA"),"Microsoft","Windows","Start Menu","Programs")
 desktop:=filepath.Join(os.Getenv("USERPROFILE"),"Desktop")
 if os.MkdirAll(start,0700)==nil{_ = shortcut(filepath.Join(start,"NOVA OS.lnk"),dest)}
 if _,e:=os.Stat(desktop);e==nil{_ = shortcut(filepath.Join(desktop,"NOVA OS.lnk"),dest)}
 setText(progressLabel,"Kurulum tamamlandı.")
 setText(statusLabel,"NOVA OS açılıyor...")
 cmd:=exec.Command(dest,"--portable");cmd.SysProcAttr=&syscall.SysProcAttr{HideWindow:true};return cmd.Start()
}
func installOffline()error{
 if len(fallbackPackage)==0{return errors.New("Çevrimdışı paket bu sürümde yok")}
 temp,e:=os.MkdirTemp("","nova-offline-*");if e!=nil{return e};defer os.RemoveAll(temp)
 path:=filepath.Join(temp,"fallback.zip")
 if e=os.WriteFile(path,fallbackPackage,0600);e!=nil{return e}
 setText(statusLabel,"İndirilemedi. EXE içindeki yedek sürüm kuruluyor...")
 return installZip(path)
}
func notifyChoice(s string,kind uintptr)uintptr{r,_,_:=procMessageBox.Call(mainWindow,uintptr(unsafe.Pointer(w(s))),uintptr(unsafe.Pointer(w("NOVA OS Kurulum"))),kind);return r}
func doInstall(){
 // Opening the hosting page also lets the user inspect the downloaded file's source.
 showSite()
 ctx,cancel:=context.WithTimeout(context.Background(),25*time.Minute);defer cancel()
 err:=install(ctx)
 if err!=nil {
  setText(statusLabel,"GitHub'dan indirme başarısız.")
  setText(progressLabel,"Yedek çevrimdışı sürüm kullanılabilir.")
  choice:=notifyChoice("GitHub paketi indirilemedi:\n\n"+err.Error()+"\n\nEXE'nin içindeki NOVA OS yedek sürümünü şimdi kurmak ister misin?",0x4|0x30)
  if choice==6 {err=installOffline()}
 }
 if err!=nil {setText(statusLabel,"Kurulum tamamlanamadı.");notify("NOVA OS kurulum hatası:\n\n"+err.Error(),MB_OK|MB_ICONERROR)}else{notify("NOVA OS başarıyla kuruldu ve başlatıldı.\n\nSanal dosyaların korunur.",MB_OK|MB_ICONINFORMATION)}
 busy=false;procEnableWindow.Call(installButton,1)
}
func main(){runtime.LockOSThread();if err:=createWindow();err!=nil{fmt.Fprintln(os.Stderr,err)}}