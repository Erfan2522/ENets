# ENets by Erfan

**A lightweight Android hotspot manager powered by Magisk and VirtualAP.**

**مدیریت‌کننده سبک هات‌اسپات اندروید با استفاده از Magisk و VirtualAP**

[English](#english) | [فارسی](#فارسی)

---

<a name="english"></a>
## English

### About

ENets is an open-source Magisk module for managing Wi-Fi hotspots on rooted Android devices through a lightweight local web dashboard.

The project uses a customized, pinned version of [VirtualAP](https://github.com/ravindu644/VirtualAP) as its hotspot backend.

> **Status: v0.2.0-alpha integration preview.** This is an experimental build. Device and firmware compatibility must be evaluated individually.

### Features

- Start and stop the Wi-Fi access point.
- Configure the hotspot name (SSID) and password.
- Select Wi-Fi security mode.
- Select 2.4 GHz or 5 GHz.
- Configure channel and channel width.
- Hide the SSID.
- View hotspot status and a read-only connected-client list.
- Responsive local web dashboard.
- Magisk-based backend service.

The dashboard is served only on `127.0.0.1:8080`. Do not expose it to the LAN or hotspot clients.

### Installation

1. Download the ENets ARM64 Magisk module ZIP from the project's [Actions](https://github.com/Erfan2522/ENets/actions/workflows/build-magisk.yml) or Releases page.
2. Open the Magisk application.
3. Go to **Modules** and select **Install from storage**.
4. Choose the downloaded ZIP.
5. Reboot your device.
6. Open this address in your device browser: `http://127.0.0.1:8080`

Keep a way to disable or remove the Magisk module if the device behaves unexpectedly.

### Compatibility

- Root access and Magisk are required.
- Current preview package: ARM64 (`arm64-v8a`).
- Android 8.0+ is an intended target, not a verified minimum.
- Compatibility depends on the Wi-Fi driver, AP interface, regulatory domain, kernel networking features, and firmware.
- The backend currently expects a usable physical Wi-Fi interface and creates/uses `ap0`; other interface layouts may need a device-specific adapter.
- A successful engine start does not by itself prove that clients can connect or access the internet.

**User-reported test device:**
- POCO M5s
- Android 17
- Magisk 30.7

The user reported successful installation, dashboard access, hotspot startup, client connection, and internet connectivity on this device. These results are not a guarantee of compatibility with other devices.

### Current limitations

- No 6 GHz support.
- No client blocking.
- No firewall configuration.
- No traffic accounting.
- No VPN routing.
- No maximum-client limit.

### Backend and upstream attribution

ENets is based on [VirtualAP](https://github.com/ravindu644/VirtualAP), pinned at commit:

`19718afbf5df48f12d82a8a30bf205cf7bfa6c6e`

The integration includes ENets-specific patches for its runtime data directory and regulatory-domain handling. See [UPSTREAM_VIRTUALAP.md](docs/UPSTREAM_VIRTUALAP.md).

### Documentation

- [Installation](docs/INSTALLATION.md)
- [Architecture](docs/ARCHITECTURE.md)
- [Compatibility](docs/COMPATIBILITY.md)
- [Roadmap](docs/ROADMAP.md)

### License

ENets is licensed under **GNU GPL-3.0-only**. See [LICENSE](LICENSE).

---

<a name="فارسی"></a>
## فارسی

### معرفی

ENets یک پروژه متن‌باز برای مدیریت هات‌اسپات وای‌فای در دستگاه‌های اندرویدی روت‌شده است.

این پروژه با استفاده از ماژول Magisk و موتور VirtualAP، یک داشبورد تحت وب محلی در اختیار کاربر قرار می‌دهد تا بتواند تنظیمات هات‌اسپات را مدیریت کند.

**وضعیت فعلی:** نسخه آزمایشی `v0.2.0-alpha`  
**شاخه توسعه:** `integration/virtualap-backend`  
**مجوز:** GNU GPL-3.0-only

### امکانات

- روشن و خاموش کردن هات‌اسپات
- تعیین نام شبکه (SSID)
- تنظیم رمز عبور
- انتخاب نوع امنیت شبکه
- انتخاب باند ۲.۴ یا ۵ گیگاهرتز
- تنظیم کانال و پهنای کانال
- مخفی کردن نام شبکه
- مشاهده وضعیت هات‌اسپات
- مشاهده فهرست دستگاه‌های متصل
- داشبورد تحت وب واکنش‌گرا
- اجرای سرویس موردنیاز از طریق Magisk

داشبورد فقط روی آدرس محلی `127.0.0.1:8080` اجرا می‌شود. آن را در شبکه محلی یا برای کاربران متصل به هات‌اسپات قابل دسترسی نکنید.

### روش نصب

۱. فایل ZIP ماژول ENets مخصوص ARM64 را از بخش [Actions](https://github.com/Erfan2522/ENets/actions/workflows/build-magisk.yml) یا Releases دریافت کنید.

۲. برنامه Magisk را باز کنید.

۳. وارد بخش **Modules** شوید و گزینه **Install from storage** را انتخاب کنید.

۴. فایل ZIP دانلودشده را انتخاب کنید.

۵. پس از نصب، دستگاه را راه‌اندازی مجدد کنید.

۶. مرورگر گوشی را باز کرده و آدرس زیر را وارد کنید:

`http://127.0.0.1:8080`

اگر دستگاه رفتار غیرعادی داشت، امکان غیرفعال‌سازی یا حذف ماژول Magisk را در نظر داشته باشید.

### سازگاری

- دسترسی روت و نصب Magisk الزامی است.
- بسته آزمایشی فعلی برای معماری ARM64 (`arm64-v8a`) است.
- اندروید ۸ به بالا هدف توسعه است، اما حداقل نسخه تأییدشده محسوب نمی‌شود.
- سازگاری نهایی به درایور وای‌فای، رابط اکسس‌پوینت، محدودیت‌های منطقه‌ای، کرنل و فریمور دستگاه بستگی دارد.
- موتور فعلی به یک رابط فیزیکی وای‌فای قابل استفاده نیاز دارد و از رابط `ap0` استفاده می‌کند؛ ممکن است دستگاه‌های دیگر به سازگارساز اختصاصی نیاز داشته باشند.
- روشن شدن موتور به‌تنهایی به معنی اتصال موفق کاربران یا دسترسی آن‌ها به اینترنت نیست.

**دستگاه آزمایش‌شده بر اساس گزارش کاربر:**
- POCO M5s
- Android 17
- Magisk 30.7

طبق گزارش کاربر، نصب برنامه، اجرای داشبورد، روشن شدن هات‌اسپات، اتصال دستگاه دیگر و دسترسی به اینترنت با موفقیت انجام شده است. این نتیجه تضمین‌کننده سازگاری با سایر دستگاه‌ها نیست.

### محدودیت‌های فعلی

- پشتیبانی از باند ۶ گیگاهرتز وجود ندارد.
- مسدودسازی دستگاه‌های متصل پیاده‌سازی نشده است.
- مدیریت فایروال وجود ندارد.
- نمایش مصرف ترافیک وجود ندارد.
- مسیریابی VPN وجود ندارد.
- محدودیت تعداد کاربران قابل تنظیم نیست.

### موتور اصلی و اعتبار پروژه

ENets بر پایه پروژه متن‌باز [VirtualAP](https://github.com/ravindu644/VirtualAP) توسعه داده شده است.

شناسه نسخه استفاده‌شده:

`19718afbf5df48f12d82a8a30bf205cf7bfa6c6e`

تغییرات اختصاصی ENets شامل مسیر داده مستقل و اصلاح نحوه مدیریت تنظیمات منطقه‌ای وای‌فای است. جزئیات بیشتر در [مستندات VirtualAP](docs/UPSTREAM_VIRTUALAP.md) قرار دارد.

### مستندات

- [راهنمای نصب](docs/INSTALLATION.md)
- [معماری پروژه](docs/ARCHITECTURE.md)
- [سازگاری دستگاه‌ها](docs/COMPATIBILITY.md)
- [نقشه راه توسعه](docs/ROADMAP.md)

### مجوز

ENets تحت مجوز **GNU GPL نسخه ۳، فقط همین نسخه (GPL-3.0-only)** منتشر می‌شود.

متن کامل مجوز در فایل [LICENSE](LICENSE) قرار دارد.
