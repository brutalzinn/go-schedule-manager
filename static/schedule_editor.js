document.getElementById('uploadForm').addEventListener('submit', async function (e) {
    e.preventDefault();
    const formData = new FormData();
    formData.append('csvfile', document.getElementById('csvfile').files[0]);

    const response = await fetch('/upload', {
        method: 'POST',
        body: formData
    });

    const result = await response.text();
    console.log('Upload result:', result);
});

document.getElementById('scheduleForm').addEventListener('submit', async function (e) {
    e.preventDefault();

    const formData = new FormData();
    formData.append('content', document.getElementById('content').value);
    formData.append('time', document.getElementById('time').value);
    formData.append('useTTS', document.getElementById('useTTS').checked);

    const response = await fetch('/create', {
        method: 'POST',
        body: formData
    });

    const result = await response.text();
    console.log('Create schedule result:', result);
});
