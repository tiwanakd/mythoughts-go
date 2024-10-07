//add button feature to show and hide the textarea
const newBtn = document.getElementById('new-btn');
const closeBtn = document.getElementById('close-btn');
const formContainer = document.querySelector('.new-tg-containter');
const btnContainer = document.querySelector('.new-tg-btn-container')

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