document.addEventListener("DOMContentLoaded", function() {
    setActiveCategory();
})

//sets active category in top bars
function setActiveCategory() {
    let cat = window.location.href.split("/").pop();
    document.querySelector(`#${cat}Btn`).setAttribute('class','active');
}
