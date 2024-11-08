//Empty out the content textarea on each page refresh;w
document.addEventListener("DOMContentLoaded", function() {
    if (textArea != null){
        textArea.value = "";
    }
    if (contentErrordiv != null) {
        contentErrordiv.style.display = "none"
    }


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
const newThoughtNavLink = document.getElementById("newThoughtNavLink")

//Show the form and hide the New button when the New button is clicked
if (newBtn != null) {
    newBtn.addEventListener("click", () => {
        formContainer.style.display = 'block';
        closeBtn.style.display = 'inline-block'
        newBtn.style.display = 'none'
        btnContainer.classList.add('close-far-end')
    });
}

if (newThoughtNavLink!= null) {
    newThoughtNavLink.addEventListener("click", () => {
        sessionStorage.setItem('showForm', 'true');
        window.location.href = "/";
    });
}

window.addEventListener("DOMContentLoaded", () => {
    if (sessionStorage.getItem('showForm') === 'true') {
        formContainer.style.display = 'block';
        closeBtn.style.display = 'inline-block';
        newBtn.style.display = 'none';
        btnContainer.classList.add('close-far-end');
        sessionStorage.removeItem('showForm');
    }
});

//when the close button is clicked hide the form and close button
if (closeBtn != null){
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
}

if (bthThoughtClear != null){
    bthThoughtClear.addEventListener("click", () => {
        textArea.value = "";
        const contentErrorlbl = document.getElementById("content-error-lbl") 
        if (contentErrorlbl != null) {
            contentErrorlbl.style.display = 'none';
            contentErrorlbl.value = '';
            textArea.classList.remove("error-field");
        }
    });
}

if (btnThoughtPost != null){
    btnThoughtPost.addEventListener("click", () => {
        const contentErrorlbl = document.getElementById("content-error-lbl");
        if (contentErrorlbl != null){
            contentErrorlbl.style.display = 'inline'
        }
    })
}

const editNameBtn = document.getElementById("edit-name")
const editUsernameBtn = document.getElementById("edit-username")
const editEmailBtn = document.getElementById("edit-email")

const editNameInput = document.getElementById("edit-name-input")
var nameValue = ""; 
if (editNameInput != null){ 
    nameValue = editNameInput.value;
}

const editUsernameInput = document.getElementById("edit-username-input")
var usernameValue = "";
if (editUsernameInput != null){
    usernameValue = editUsernameInput.value;
}

const editEmailInput = document.getElementById("edit-email-input")
var emailValue = "";
if (editEmailInput != null){
    emailValue = editEmailInput.value;
}

const editNameSubmit = document.getElementById("edit-name-submit")
const editUsernameSubmit = document.getElementById("edit-username-submit")
const editEmailSubmit = document.getElementById("edit-email-submit")
const editNameCancel = document.getElementById("edit-name-cancel")
const editUsernameCancel = document.getElementById("edit-username-cancel")
const editEmailCancel = document.getElementById("edit-email-cancel")

const userEditNameError = document.getElementById("content-error-name")
const userEditUsernameError = document.getElementById("content-error-username")
const userEditEmailError = document.getElementById("content-error-email")

const contentErrorlbl = document.getElementById("content-error-lbl");
const userNameErrorlbl = document.getElementById("user-name-error-lbl")
const userUameErrorlbl = document.getElementById("user-username-error-lbl")
const userEmailErrorlbl = document.getElementById("user-email-error-lbl")

const currentPasswordField = document.getElementById("curr-pass")
const newPasswordField = document.getElementById("new-pass")
const confirmPasswordField = document.getElementById("conf-pass")

document.addEventListener("htmx:afterOnLoad", function(event) {
    // Check for 422 status to apply error styles and keep form visible
    if (event.detail.xhr.status === 422) {
        if (contentErrordiv != null){
            contentErrordiv.style.display = 'inline-block'
        }
        if (textArea != null) {
            textArea.classList.add("error-field");
        }

        const button = event.detail.elt; 
        const url = button ? button.getAttribute("hx-put") : null;
        
        if (url == "/user/account/edit/name"){
            if (userEditNameError != null){
                userEditNameError.style.display = 'block'
            }

            if (editNameInput != null) {
                editNameInput.classList.add("error-field")
            }
        }

        if (url == "/user/account/edit/username"){
            if (userEditUsernameError != null){
                userEditUsernameError.style.display = 'block'
            }

            if (editUsernameInput != null) {
                editUsernameInput.classList.add("error-field")
            }
        }

        if (url == "/user/account/edit/email"){
            if (userEditEmailError != null){
                userEditEmailError.style.display = 'block'
            }

            if (editEmailInput != null) {
                editEmailInput.classList.add("error-field")
            }
        }

        if (url == "/user/account/password/update"){
            if (currentPasswordField != null  && newPasswordField != null) {
                currentPasswordField.value = '';
                newPasswordField.value = '';
                confirmPasswordField.value = '';
                currentPasswordField.classList.add("error-field")
                newPasswordField.classList.add("error-field")
                confirmPasswordField.classList.add("error-field")
            }
        }

    } else if (event.detail.xhr.status === 201) {
        // Hide form on successful POST (200)
        if (textArea != null){
            textArea.value = "";
            textArea.classList.remove("error-field");
        }
        if (contentErrorlbl != null){
            contentErrorlbl.style.display = 'none'
            contentErrorlbl.value = ''
        }

        const thoughtFormCont = document.getElementById("thought-form-container")
        if (thoughtFormCont != null){
            thoughtFormCont.style.display = 'none';
        }
        if (newBtn != null){
            newBtn.style.display = 'inline';
        }

        if (closeBtn != null){
            closeBtn.style.display = 'none';
        }

    }
});

clearUserFormBtn = document.getElementById("clear-user-form");
if (clearUserFormBtn != null){
        clearUserFormBtn.addEventListener("click", () => {
        document.getElementById("name-input").value = ""; 
        document.getElementById("username-input").value = ""; 
        document.getElementById("email-input").value = ""; 
        document.getElementById("password-input").value = ""; 
    });
}

var navLinks = document.querySelectorAll("nav a");
for (var i = 0; i < navLinks.length; i++) {
	var link = navLinks[i]
	if (link.getAttribute('href') == window.location.pathname) {
		link.classList.add("live");
		break;
	}
}

if (editNameBtn != null) {
    editNameBtn.addEventListener("click", () => {
        editNameInput.classList.remove('hidden');
        editNameSubmit.classList.remove('hidden');
        editNameSubmit.classList.add('edit-submit-btn')
        editNameCancel.classList.remove('hidden');
        editNameCancel.classList.add('edit-cancel-btn');
    });
}

if (editUsernameBtn != null) {
    editUsernameBtn.addEventListener("click", () => {
        editUsernameInput.classList.remove('hidden');
        editUsernameSubmit.classList.remove('hidden');
        editUsernameSubmit.classList.add('edit-submit-btn')
        editUsernameCancel.classList.remove('hidden');
        editUsernameCancel.classList.add('edit-cancel-btn');
    })
}

if (editEmailBtn != null) {
    editEmailBtn.addEventListener("click", () => {
        editEmailInput.classList.remove('hidden');
        editEmailSubmit.classList.remove('hidden');
        editEmailSubmit.classList.add('edit-submit-btn')
        editEmailCancel.classList.remove('hidden');
        editEmailCancel.classList.add('edit-cancel-btn');
    })
}

if (editNameCancel != null) {
    editNameCancel.addEventListener("click", () => {
        editNameInput.classList.add('hidden');
        editNameSubmit.classList.add('hidden');
        editNameCancel.classList.add('hidden');
        editNameInput.classList.remove('error-field')
        editNameInput.value = nameValue;
        if (userNameErrorlbl != null){
            userNameErrorlbl.value = ''
        }
        userEditNameError.style.display = 'none'
    })
}

if (editUsernameCancel != null) {
    editUsernameCancel.addEventListener("click", () => {
        editUsernameInput.classList.add('hidden');
        editUsernameSubmit.classList.add('hidden');
        editUsernameCancel.classList.add('hidden');
        editUsernameInput.classList.remove('error-field')
        editUsernameInput.value = usernameValue;
        if (userUameErrorlbl != null){
            userUameErrorlbl.value = ''       
        }
        userEditUsernameError.style.display = 'none'
    })
}

if (editEmailCancel != null) {
    editEmailCancel.addEventListener("click", () => {
        editEmailInput.classList.add('hidden');
        editEmailSubmit.classList.add('hidden');
        editEmailCancel.classList.add('hidden');
        editEmailInput.classList.remove('error-field')
        editEmailInput.value = emailValue;
        if (userEmailErrorlbl != null){
            userEmailErrorlbl.value = ''
        }
        userEditEmailError.style.display = 'none'
    })
}

