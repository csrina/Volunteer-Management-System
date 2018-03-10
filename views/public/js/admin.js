document.addEventListener("DOMContentLoaded", function() {
    setActiveCategory();
    switch (window.location.href.split("/").pop()) {
        case 'dashboard':
            loadDash();
            break;
    }
})

//sets active category in top bars
function setActiveCategory() {
    let cat = window.location.href.split("/").pop();
    document.querySelector(`#${cat}Btn`).setAttribute('class','active');
}

function loadDash() {
}
