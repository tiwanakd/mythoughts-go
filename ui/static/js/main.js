//Implement the like/dislike feature
document.addEventListener("DOMContentLoaded", () => {
    const likeForms = document.querySelectorAll('[id^="like-form-"]');
    const dislikeForms = document.querySelectorAll('[id^="dislike-form-"]');

    likeForms.forEach(form => form.addEventListener("submit", handleLikeSubmit))
    dislikeForms.forEach(form => form.addEventListener("submit", handleDislikeSubmit))
})

function handleLikeSubmit(event){
    event.preventDefault()
    const id = event.target.id.split('-')[2];
    addLike(id)
}

function handleDislikeSubmit(event){
    event.preventDefault()
    const id = event.target.id.split('-')[2];
    addDislike(id)
}

function addLike(id){
    const url = `/like/${id}`;
    fetch(url, {
        method: 'POST'
    })
    .then(response => response.json())
    .then(data => {
        const agreeCountElement = document.getElementById(`agree-count-${id}`);
        agreeCountElement.textContent = data.newAgreeCount;
    })
    .catch(error => {
        console.error('Error', error)
    })

}

function addDislike(id){
    const url = `dislike/${id}`;
    fetch(url, {
        method: 'POST'
    })
    .then(resp => resp.json())
    .then(data => {
        const disgreeCountElement = document.getElementById(`disagree-count-${id}`);
        disgreeCountElement.textContent = data.newDisagreeCount;
    })
    .catch(error => {
        console.error('Error', error)
    })
}

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