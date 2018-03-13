function check() {
    var new_pwd_txt = document.getElementById('new_pwd');
    var conf_pwd_txt =  document.getElementById('confirm_pwd');
    var conf_pwd_lbl = document.getElementById('confirm_pwd_lbl');
    var message_small = document.getElementById('message');
    console.log(new_pwd_txt.value);
    console.log(conf_pwd_lbl.value);
    if (new_pwd_txt.value == conf_pwd_txt.value) {
        message_small.classList.add("d-none");
        conf_pwd_lbl.classList.remove("text-danger");
        conf_pwd_lbl.classList.add("text-success");
        conf_pwd_txt.classList.remove("is-invalid");
        conf_pwd_txt.classList.add("is-valid");

    } else if (new_pwd_txt.value != conf_pwd_txt.value && conf_pwd_txt.value.length > 0){
        message_small.classList.remove("d-none");
        conf_pwd_lbl.classList.add("text-danger");
        conf_pwd_lbl.classList.remove("text-success");
        conf_pwd_txt.classList.add("is-invalid");
        conf_pwd_txt.classList.remove("is-valid");
    } else if (conf_pwd_txt.value.length == 0) {
        message_small.classList.add("d-none");
        conf_pwd_lbl.classList.remove("text-danger");
        conf_pwd_lbl.classList.remove("text-success");
        conf_pwd_txt.classList.remove("is-invalid");
        conf_pwd_txt.classList.remove("is-valid");  
    }
}
function change_pwd() {

    alert("Change password is not yet implemented")
}