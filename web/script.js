document.getElementById('feedbackForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    
    const status = document.getElementById('status');
    const submitBtn = document.getElementById('submitBtn');
    
    // ۱. استخراج داده‌ها (شامل فیلد مخفی Honeypot)
    const message = document.getElementById('message').value;
    const hp = document.getElementById('hp').value;

    status.innerText = "در حال ارسال امن...";
    status.style.color = "black";
    submitBtn.disabled = true;

    try {
        // ۲. ارسال به Vercel (آدرس را بعد از Deploy جایگزین کن)
        const response = await fetch('https://whisper-bridge-api.vercel.app/api/submit', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ message, hp })
        });

        // ۳. مدیریت پاسخ‌های مختلف بک‌ایند
        if (response.status === 201) {
            status.style.color = "green";
            status.innerText = "✓ پیام شما بدون ذخیره هیچ ردپایی، با موفقیت بایگانی شد.";
            document.getElementById('feedbackForm').reset();
        } else if (response.status === 429) {
            status.style.color = "orange";
            status.innerText = "⚠️ شما در این دقیقه بیش از حد پیام فرستاده‌اید. کمی صبر کنید.";
        } else {
            throw new Error();
        }
    } catch (error) {
        status.style.color = "red";
        status.innerText = "❌ خطا در اتصال. لطفاً وضعیت اینترنت خود را چک کنید.";
    } finally {
        submitBtn.disabled = false;
    }
});
