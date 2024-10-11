//Empty out the content textarea on each page refresh;w
document.addEventListener("DOMContentLoaded", function() {
    textArea.value = "";
    contentErrordiv.style.display = "none"
});

//add button feature to show and hide the textarea
const newBtn = document.getElementById('new-btn');
const closeBtn = document.getElementById('close-btn');
const formContainer = document.querySelector('.new-tg-containter');
const btnContainer = document.querySelector('.new-tg-btn-container');
const bthThoughtClear = document.getElementById("thought-clear-btn");
const btnThoughtPost = document.getElementById("thought-post-btn");
const textArea = document.getElementById("new-tg-box");
const contentErrordiv = document.getElementById("content-error")

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
    contentErrordiv.style.display = 'none';
    const contentErrorlbl = document.getElementById("content-error-lbl") 
    if (contentErrorlbl != null) {
        contentErrorlbl.style.display = 'none';
        contentErrorlbl.value = '';
        textArea.classList.remove("error-field");
    }
});

bthThoughtClear.addEventListener("click", () => {
    textArea.value = "";
    const contentErrorlbl = document.getElementById("content-error-lbl") 
    if (contentErrorlbl != null) {
        contentErrorlbl.style.display = 'none';
        contentErrorlbl.value = '';
        textArea.classList.remove("error-field");
    }
})

btnThoughtPost.addEventListener("click", () => {
    const contentErrorlbl = document.getElementById("content-error-lbl");
    if (contentErrorlbl != null){
        contentErrorlbl.style.display = 'inline'
    }
})

document.addEventListener("htmx:afterOnLoad", function(event) {
    // Check for 422 status to apply error styles and keep form visible
    if (event.detail.xhr.status === 422) {
        contentErrordiv.style.display = 'inline-block'
        textArea.classList.add("error-field");
    } else if (event.detail.xhr.status === 201) {
        // Hide form on successful POST (200)
        textArea.value = "";
        const contentErrorlbl = document.getElementById("content-error-lbl");
        if (contentErrorlbl != null){
            contentErrorlbl.style.display = 'none'
            contentErrorlbl.value = ''
        }
        document.getElementById("thought-form-container").style.display = 'none';
        newBtn.style.display = 'inline';
        closeBtn.style.display = 'none';
        textArea.classList.remove("error-field"); // Remove error styles on success
    }
});