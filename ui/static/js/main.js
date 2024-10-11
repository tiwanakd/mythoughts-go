//Empty out the content textarea on each page refresh;w
document.addEventListener("DOMContentLoaded", function() {
    textArea.value = "";
});

//add button feature to show and hide the textarea
const newBtn = document.getElementById('new-btn');
const closeBtn = document.getElementById('close-btn');
const formContainer = document.querySelector('.new-tg-containter');
const btnContainer = document.querySelector('.new-tg-btn-container');
const bthThoughtClear = document.getElementById("thought-clear-btn");
const btnThoughtPost = document.getElementById("thought-post-btn");
const textArea = document.getElementById("new-tg-box");

//Show the form and hide the New button when the New button is clicked
newBtn.addEventListener("click", () => {
    formContainer.style.display = 'block';
    closeBtn.style.display = 'inline-block'
    newBtn.style.display = 'none'
    btnContainer.classList.add('close-far-end')
});

//when the close button is clicked hide the form and close button
closeBtn.addEventListener("click", () => {
    formContainer.style.display = 'none';
    closeBtn.style.display = 'none';
    newBtn.style.display = 'inline-block';
    btnContainer.classList.remove('close-far-end');
    document.getElementById("content-error-lbl").style.display = 'none';
});

bthThoughtClear.addEventListener("click", () => {
    textArea.value = "";
})

document.addEventListener("htmx:afterOnLoad", function(event) {
    const contentErrorLabel = document.getElementById("content-error-lbl")
    // Check for 422 status to apply error styles and keep form visible
    if (event.detail.xhr.status === 422) {
        textArea.classList.add("error-field");
    } else if (event.detail.xhr.status === 200) {
        // Hide form on successful POST (200)
        textArea.value = "";
        contentErrorLabel.style.display = 'none';
        contentErrorLabel.value = '';
        document.getElementById("thought-form-container").style.display = 'none';
        newBtn.style.display = 'inline';
        closeBtn.style.display = 'none';
        textArea.classList.remove("error-field"); // Remove error styles on success
    }
});