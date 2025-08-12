console.log('contact.js loaded');

// Initialize EmailJS with user ID
    emailjs.init('-eI6FqIyAwQ91X7-J'); // USER_ID is now the public key on emailJS
    const sendButton = document.getElementById('send-button');
    const buttonText = document.getElementById('button-text');
    const buttonSpinner = document.getElementById('button-spinner');

    document.getElementById('email-form').addEventListener('submit', function (e) {
        e.preventDefault();

        // Clear previous errors and status
        ['name', 'email', 'message'].forEach(field => {
            document.getElementById(`error-${field}`).textContent = '';
        });
        const statusEl = document.getElementById('form-status');
        statusEl.textContent = '';
        statusEl.className = '';

        const name = document.getElementById('sender-name').value.trim();
        const email = document.getElementById('sender-email').value.trim();
        const message = document.getElementById('message').value.trim();

        let valid = true;

        if (!name) {
            document.getElementById('error-name').textContent = 'Please enter your name';
            valid = false;
        }
        if (!email) {
            document.getElementById('error-email').textContent = 'Please enter your email';
            valid = false;
        } else if (!validateEmail(email)) {
            document.getElementById('error-email').textContent = 'Please enter a valid email';
            valid = false;
        }
        if (!message) {
            document.getElementById('error-message').textContent = 'Please enter a message';
            valid = false;
        }
        if (!valid) return;

        // Show spinner & hide button text
        buttonText.classList.add('hidden');
        buttonSpinner.classList.remove('hidden');
        sendButton.disabled = true;
        sendButton.classList.add('opacity-50', 'cursor-not-allowed');


        // Prepare EmailJS template parameters
        const templateParams = {
            from_name: name,
            from_email: email,
            message: message,
        };

        emailjs.send('service_xze3q5s', 'template_jv56516', templateParams) // Service ID and Template ID are safe to expose since they're basically just public keys
            .then(() => {
                buttonText.classList.remove('hidden');
                buttonSpinner.classList.add('hidden');
                sendButton.disabled = false;
                sendButton.classList.remove('opacity-50', 'cursor-not-allowed');

                statusEl.textContent = 'Sent!';
                statusEl.className = 'text-green-600 text-center';
                document.getElementById('email-form').reset();

                // Remove the "Sent!" message after 3 seconds
                setTimeout(() => {
                    statusEl.textContent = '';
                    statusEl.className = '';
                }, 3000);
            }, (error) => {
                buttonText.classList.remove('hidden');
                buttonSpinner.classList.add('hidden');
                sendButton.disabled = false;
                sendButton.classList.remove('opacity-50', 'cursor-not-allowed');
            
                statusEl.textContent = 'Failed to send message. Please try again later.';
                statusEl.className = 'text-red-600';
                console.error('EmailJS error:', error);
            });
    });

    function validateEmail(email) {
        return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);
    }