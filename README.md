# Bilgi Arena

2–4 kişilik Türkçe WebSocket bilgi yarışması. Node.js 22 veya üzeri gerekir.

## Çalıştırma

npm install ve npm start komutlarını çalıştırın. Varsayılan port 3000; PORT ortam değişkeniyle değişir. npm test dört oyunculu oda ve puan güvenliği testlerini çalıştırır.

## Özellikler

6 karakterli oda kodu, davet bağlantısı, 4 oyuncu sınırı, yalnız oda sahibinin başlatması, 10 soru, soru başına 15 saniye, sunucuda cevap doğrulama ve süre hesaplama. Doğru cevap 1000 puan ve en çok 500 hız puanı verir. Cevap ve puan artışı tur bitene kadar açıklanmaz. Cevap değiştirilemez. 60 saniye yeniden bağlanma, sahiplik aktarımı, beraberlik ve yeniden oynama vardır.

320 soru; 8 kategori ve Karışık seçeneği; Kolay, Orta, Zor ve Karışık zorluk. Türkçe koyu arayüz, avatarlar, sıralama, ses anahtarı, konfeti ve mobil yerleşim.

## Yayınlama

Klasörü GitHub deponuzun köküne yükleyin. Render üzerinde Node Web Service: build npm install --omit=dev --ignore-scripts; start node server.mjs; Free plan. Alternatif olarak render.yaml dosyasından Blueprint oluşturulabilir. Dockerfile dahildir.

Tek sunucu örneği kullanın. Odalar bellektedir; yeniden başlatmada kaybolur. Ücretsiz Render hareketsizlikten sonra uyuyabilir. Güncel koşullar: https://render.com/docs/free

## Doğrulama ve sınırlar

30 Ekim 2026: dört WebSocket oyuncusu, dolu oda reddi, oda sahibi yetkisi, aynı soru/süre, erken cevap gizliliği, 10 sorunun bitmesi, sıfır cevap/sıfır puan, yeniden oynama, sahiplik aktarımı ve yeniden bağlanma testleri geçti. Doğru/yanlış puanlama, tekrar cevap ve eski soruya cevap reddi geçti. Önceki sürüm gerçek HTTPS üzerinden iki tarayıcı sekmesinde oynandı.

Farklı fiziksel telefon testi ve son mobil görsel kontrol yapılmadı. Kalıcı Render yayını tamamlanmadı. Geçici Cloudflare bağlantısı yerel sunucu ve tünel açıkken çalışır; 7/24 barındırma değildir. Soruların yalnızca 19'u bu çalışmada resmî kaynakla tek tek karşılaştırıldı; SORU-KAYNAKLARI.md dosyasına bakın.
