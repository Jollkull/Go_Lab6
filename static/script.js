document.addEventListener('DOMContentLoaded', () => {
    const form = document.querySelector('form');
    const result = document.getElementById('result');

    if (form) {
        form.addEventListener('submit', async (e) => {
            e.preventDefault();

            const inputs = [...form.querySelectorAll('input')].map(input => {
                const value = parseFloat(input.value);
                return isNaN(value) ? 0 : value;
            });

            const url = form.id === 'calculator1' ? '/api/calculator1' : '/api/calculator2';

            try {
                result.textContent = "Calculating...";
                
                const response = await fetch(url, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ values: inputs })
                });

                const data = await response.json();
                result.textContent = `Result: ${data.result}`;
                result.style.color = '#2ecc71';
            } catch (error) {
                result.textContent = 'Error calculating. Please try again.';
                result.style.color = '#e74c3c';
                console.error('Calculation error:', error);
            }
        });
    }
});