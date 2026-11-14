# NOVA OS v0.2 Beta

**Windows 11 tarzı, Windows içinde çalışan sanal masaüstü.** VirtualBox veya Linux değildir; Windows donanımını değiştirmez ve bağımsız işletim sistemi değildir. Kullanıcının Windows programlarını sanal Windows içine kurmasına izin vermez.

## Kurulum modeli

- `NOVA_OS_WebInstaller.exe` (~6 MB; içinde çevrimdışı yedek de bulunur) internet indirme başlatıcısıdır.
- `KUR VE İNDİR`: GitHub Releases sitesini açar, `SHA256SUMS.txt` ve `NOVA_OS_Package.zip` indirir, paket bütünlüğünü kontrol eder, EXE'yi `%LOCALAPPDATA%\Programs\NOVA OS` içine kurup başlatır.
- Kullanıcı verileri `%LOCALAPPDATA%\NOVA OS\Profile` altında kalır. Güncelleme önceki profil klasörünü silmez.
- Yönetici yetkisi ve ayrı kurulum yazılımı gerekmez; Windows 11 ve Edge gerekir.
- GitHub sitesi bulunamazsa veya sürüm yayımlanmamışsa bir hata gösterilir ve **EXE içine gömülü yedek sürümü kurmak için izin istenir**. Bu yedek sürüm internet olmadan da çalışabilir. Alternatif olarak `NOVA_OS.exe` indirilebilir.

Yayın sayfası: https://github.com/ugurfc26-eng/ugur/releases/tag/nova-os-v0.2.0

## Uygulamalar

Dosyalar, notlar, metin editörü, hesap makinesi, web bağlantıları, çizim, fotoğraf galerisi, ses/video oynatma, sanal terminal, ajanda, saat/kronometre, NOVA arama, görev yöneticisi, ayarlar, sistem hakkında. Dosyalar IndexedDB'de yerel ve kalıcı saklanır. TXT ve PNG indirme, bilgisayardan dosya içeri alma desteklenir.

## Sistem ve sınırlar

- Sanal Windows 11 kurulumu değildir ve Windows EXE programlarını ayrı bir sanal makinede çalıştıramaz.
- Uygulama RAM'i arttırmaz; ekran kartı emülatörü içermez.
- Ağ tarayıcı bağlantıları Windows'un varsayılan tarayıcısında açılır.
- SmartScreen uyarısı çıkabilir: bu bir imzasız uygulamadır.
- Verilerinizin yedeğini almak için içeri aktardığınız veya ürettiğiniz dosyaları dışarı aktarın; Edge kullanıcı verisi profilini silmeyin.
- Tam Windows 11 masaüstünde uçtan uca gerçek makine testi yapılmadı. Beta sürümüdür.

## Güvenlik

Web installer yalnızca bu GitHub deposundaki sabit sürüm paketini kullanır, HTTPS indirir, SHA256 ile bütünlük kontrolü yapar ve ZIP içinde adı `NOVA_OS.exe` olan dosyayı kurar. Kodu ve derleme süreci kaynakta açıktır; paket başka sitelerden sessizce indirilemez.