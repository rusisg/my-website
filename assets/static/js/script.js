const themeToggle = document.getElementById('theme-toggle');

// Function to apply the theme
function applyTheme(isDarkMode) {
    // Apply or remove the dark mode class
    document.body.classList.toggle('dark', isDarkMode);

    // Update SVG fills
    const svgElements = document.querySelectorAll('svg');
    svgElements.forEach(svg => {
        const paths = svg.querySelectorAll('path');
        paths.forEach(path => {
            path.setAttribute('fill', isDarkMode ? 'white' : 'black');
        });
    });

    // Update navigation button styles
    const navButtons = document.querySelectorAll('.nav-button');
    navButtons.forEach(button => {
        button.style.color = isDarkMode ? 'white' : 'black';
    });
}


// Load the saved theme preference on page load
document.addEventListener('DOMContentLoaded', () => {
    const savedTheme = localStorage.getItem('theme'); // Retrieve saved theme
    const isDarkMode = savedTheme === 'dark'; // Determine if dark mode is active
    applyTheme(isDarkMode); // Apply the theme
    themeToggle.checked = isDarkMode; // Sync toggle switch state
});

// Add event listener for theme toggle
themeToggle.addEventListener('change', () => {
    const isDarkMode = themeToggle.checked; // Check the toggle state
    applyTheme(isDarkMode); // Apply the theme
    localStorage.setItem('theme', isDarkMode ? 'dark' : 'light'); // Save preference
});


