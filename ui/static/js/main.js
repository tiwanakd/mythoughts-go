//Empty out the content textarea on each page refresh;w
document.addEventListener("DOMContentLoaded", function() {
    document.getElementById("new-tg-box").value = "";
});

//allow htmx to swap code 422
document.body.addEventListener("htmx:beforeSwap", function (evt) {
    if (evt.detail.xhr.status === 422) {
    evt.detail.shouldSwap = true;
    evt.detail.isError = false;
  } 
});

//add button feature to show and hide the textarea
const newBtn = document.getElementById('new-btn');
const closeBtn = document.getElementById('close-btn');
const formContainer = document.querySelector('.new-tg-containter');
const btnContainer = document.querySelector('.new-tg-btn-container');
const bthThoughtClear = document.getElementById("thought-clear-btn");
const btnThoughtPost = document.getElementById("thought-post-btn")

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
    closeBtn.style.display = 'none'
    newBtn.style.display = 'inline-block'
    btnContainer.classList.remove('close-far-end')
});

bthThoughtClear.addEventListener("click", () => {
    document.getElementById("new-tg-box").value = "";
})

// btnThoughtPost.addEventListener("click", () => {
//     newBtn.style.display = 'inline-block'
//     formContainer.style.display = 'none';
//     closeBtn.style.display = 'none'
// })
document.addEventListener("htmx:afterRequest", function() {
    const errorLabel = document.querySelector("#content-error .error");
    const textArea = document.getElementById("new-tg-box");

    if (errorLabel) {
        textArea.classList.add("error-field");
    } else {
        textArea.classList.remove("error-field");
    }
});