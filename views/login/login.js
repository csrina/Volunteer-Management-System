document.getElementById("loginbtn").addEventListener("click", login);

function login(){
    var name= document.querySelector("#usr").value;
    console.log("username: " + name  );

    var pass = document.querySelector("#pwd").value;
    console.log("password: " + pass  );

    var xmlhttp = new XMLHttpRequest();   // new HttpRequest instance 
    xmlhttp.open("POST", "http://localhost:8080/api/v1/login");
    xmlhttp.setRequestHeader("Content-Type", "application/json");
    xmlhttp.send(JSON.stringify({username:name, password:pass}));
}