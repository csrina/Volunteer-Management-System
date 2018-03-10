document.addEventListener("DOMContentLoaded", function() {
    setActiveCategory();
})

//sets active category in top bars
function setActiveCategory() {
    let cat = window.location.href.split("/").pop();
    document.querySelector(`#${cat}Btn`).setAttribute('class','active');
}

function userList() {
    let xhttp = new XMLHttpRequest();
    xhttp.addEventListener("loadend", () => {
        console.log(xhttp.response);
        let userInfo = JSON.parse(xhttp.response);
        let tmpl = document.querySelector("#tmpl_listUsers").innerHTML;
        let func = doT.template(tmpl);
        document.querySelector("#displayArea").innerHTML = func(userInfo);
        console.log("finished")
    });
    xhttp.open("GET", `http://localhost:8080/api/v1/admin/users`);
    xhttp.send();
}
