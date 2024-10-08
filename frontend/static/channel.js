document.getElementById('configForm').addEventListener('submit', async function(event) {
    event.preventDefault(); // 防止表單默認提交

    const channelSecret = document.getElementById('channelSecret').value;
    const channelToken = document.getElementById('channelToken').value;

    const data = {
        channelSecretKey: channelSecret,
        channelAccessToken: channelToken
    };

    try {
        const response = await fetch('/channel/create', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data)
        });

        if (response.ok) {
            const result = await response.json();
            console.log('Success:', result);
            alert('Configuration saved successfully!');
        } else {
            throw new Error('Network response was not ok');
        }
    } catch (error) {
        console.error('Error:', error);
        alert('Failed to save configuration.');
    }
});
