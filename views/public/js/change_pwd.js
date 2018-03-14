function enterPressed(e) {
    if (e.keycode == 13 || e.which == 13) {
        e.preventDefault();
        var t = check();
        if (t == true) {
            change_pwd();
        }
        else {
            console.log("not quite yet");
        }
    }
}

function check() {
    var old_pwd_txt = document.getElementById('old_pwd');
    var old_pwd_lbl = document.getElementById('old_pwd_lbl');
    var old_message = document.getElementById('old_message');

    var new_pwd_txt = document.getElementById('new_pwd');
    var new_pwd_lbl = document.getElementById('new_pwd_lbl');
    var new_message = document.getElementById('new_message');

    var conf_pwd_txt = document.getElementById('confirm_pwd');
    var conf_pwd_lbl = document.getElementById('confirm_pwd_lbl');
    var conf_message = document.getElementById('conf_message');

    var submit_btn = document.getElementById('change_pwd_btn');

    if (new_pwd_txt.value.length < 8) {
        new_message.classList.remove("d-none");
        new_pwd_lbl.classList.add("text-danger");
        new_pwd_lbl.classList.remove("text-success");
        new_pwd_txt.classList.add("is-invalid");
        new_pwd_txt.classList.remove("is-valid");
    }
    else {
        new_message.classList.add("d-none");
        new_pwd_lbl.classList.remove("text-danger");
        new_pwd_lbl.classList.add("text-success");
        new_pwd_txt.classList.remove("is-invalid");
        new_pwd_txt.classList.add("is-valid");
    }



    if (new_pwd_txt.value == conf_pwd_txt.value && new_pwd_txt.value.length != 0 && conf_pwd_txt.value.length != 0) {
        conf_message.classList.add("d-none");
        conf_pwd_lbl.classList.remove("text-danger");
        conf_pwd_lbl.classList.add("text-success");
        conf_pwd_txt.classList.remove("is-invalid");
        conf_pwd_txt.classList.add("is-valid");

    }
    else if (new_pwd_txt.value != conf_pwd_txt.value && conf_pwd_txt.value.length > 0) {
        conf_message.innerText = "Needs to match new password";
        conf_message.classList.remove("d-none");
        conf_pwd_lbl.classList.add("text-danger");
        conf_pwd_lbl.classList.remove("text-success");
        conf_pwd_txt.classList.add("is-invalid");
        conf_pwd_txt.classList.remove("is-valid");
    }
    else if (new_pwd_txt.value != conf_pwd_txt.value && conf_pwd_txt.value.length == 0) {
        conf_message.innerText = "Cannot be empty and needs to match new password";
        conf_message.classList.remove("d-none");
        conf_pwd_lbl.classList.add("text-danger");
        conf_pwd_lbl.classList.remove("text-success");
        conf_pwd_txt.classList.add("is-invalid");
        conf_pwd_txt.classList.remove("is-valid");
    }
    else if (new_pwd_txt.value.length == 0 && conf_pwd_txt.value.length == 0) {
        conf_message.innerText = "Cannot be empty";
        conf_message.classList.remove("d-none");
        conf_pwd_lbl.classList.add("text-danger");
        conf_pwd_lbl.classList.remove("text-success");
        conf_pwd_txt.classList.add("is-invalid");
        conf_pwd_txt.classList.remove("is-valid");
    }
    if (old_pwd_txt.value.length > 0) {
        old_message.classList.add("d-none");
        old_pwd_lbl.classList.remove("text-danger");
        old_pwd_txt.classList.remove("is-invalid");

    }
    else {
        old_message.classList.remove("d-none");
        old_pwd_lbl.classList.add("text-danger");
        old_pwd_txt.classList.add("is-invalid");

    }

    if (old_pwd_txt.value.length == 0 || conf_pwd_txt.value.length == 0 ||
        (new_pwd_txt.value != conf_pwd_txt.value && conf_pwd_txt.value.length > 0) ||
        new_pwd_txt.value.length < 8) {
        submit_btn.disabled = true;
        return false;
    }
    else {
        submit_btn.disabled = false;
        return true;
    }


}

function change_pwd() {
    var old_pwd_txt = document.getElementById('old_pwd');
    var old_pwd_lbl = document.getElementById('old_pwd_lbl');

    var new_pwd_txt = document.getElementById('new_pwd');
    var new_pwd_lbl = document.getElementById('new_pwd_lbl');


    var conf_pwd_txt = document.getElementById('confirm_pwd');
    var conf_pwd_lbl = document.getElementById('confirm_pwd_lbl');

    var submit_btn = document.getElementById('change_pwd_btn');


    old_pwd_lbl.classList.remove("text-danger");
    old_pwd_txt.classList.remove("is-invalid");

    old_pwd_txt.disabled = true;
    new_pwd_txt.disabled = true;
    conf_pwd_txt.disabled = true;
    submit_btn.disabled = true;

    var ready = check();
    if (ready == false) {
        old_pwd_txt.disabled = false;
        new_pwd_txt.disabled = false;
        conf_pwd_txt.disabled = false;
        submit_btn.disabled = false;
        return;
    }

    var location = document.getElementsByClassName('container');
    var progress = location[0];
    progress.insertAdjacentHTML('afterbegin',
        '<div class="progress" id="progress_div" style="margin-bottom: 20px"> <div id="progress" class="progress-bar progress-bar-striped progress-bar-animated" role="progressbar" aria-valuenow="75" aria-valuemin="0" aria-valuemax="100" style="width: 25%"></div><br></div>')


    Command: toastr["info"]("Updating Password ......")

    toastr.options = {
        "closeButton": false,
        "debug": false,
        "newestOnTop": true,
        "progressBar": false,
        "positionClass": "toast-top-full-width",
        "preventDuplicates": true,
        "onclick": null,
        "showDuration": "300",
        "hideDuration": "1000",
        "timeOut": "5000",
        "extendedTimeOut": "1000",
        "showEasing": "swing",
        "hideEasing": "linear",
        "showMethod": "fadeIn",
        "hideMethod": "fadeOut"
    }

    // check old password
    const api_call = "/api/v1/passwords";
    console.log("one");
    const xmlhttp = new XMLHttpRequest(); // new HttpRequest instance
    console.log("two");
    xmlhttp.onreadystatechange = function() {
        check();
        console.log("STUFSSSS");
        if (xmlhttp.readyState !== 4) {
            console.log("not ready");
            return;
        }
        if (xmlhttp.status !== 202) {
            console.log("Failed to auth");
            console.log(xmlhttp.status);
            console.log(xmlhttp.responseText);
            Command: toastr["error"]("Updating password failed.")

            toastr.options = {
                "closeButton": false,
                "debug": false,
                "newestOnTop": true,
                "progressBar": false,
                "positionClass": "toast-top-full-width",
                "preventDuplicates": true,
                "onclick": null,
                "showDuration": "300",
                "hideDuration": "1000",
                "timeOut": "5000",
                "extendedTimeOut": "1000",
                "showEasing": "swing",
                "hideEasing": "linear",
                "showMethod": "fadeIn",
                "hideMethod": "fadeOut"
            }
            var elem = document.getElementById('progress_div');
            elem.parentNode.removeChild(elem);

            var l = document.getElementById('alert_box');
            l.classList.remove("d-none");

            old_pwd_lbl.classList.add("text-danger");
            old_pwd_txt.classList.add("is-invalid");

            new_pwd_lbl.classList.remove("text-danger");
            new_pwd_lbl.classList.remove("text-success");
            new_pwd_txt.classList.remove("is-invalid");
            new_pwd_txt.classList.remove("is-valid");

            conf_pwd_lbl.classList.remove("text-danger");
            conf_pwd_lbl.classList.remove("text-success");
            conf_pwd_txt.classList.remove("is-invalid");
            conf_pwd_txt.classList.remove("is-valid");

            old_pwd_txt.disabled = false;
            new_pwd_txt.disabled = false;
            conf_pwd_txt.disabled = false;
            submit_btn.disabled = false;
            return;
        }
        else {
            old_pwd_txt.disabled = true;
            new_pwd_txt.disabled = true;
            conf_pwd_txt.disabled = true;
            submit_btn.disabled = true;
            var l = document.getElementById('alert_box');
            l.classList.add("d-none");
            l.innerHTML = '<strong>Warning!</strong> Password update failed, could not authenticate user. Please enter correct password.'
            var p = document.getElementById('progress')
            p.style.width = '60%';
            console.log("User authenticated");
            alert("Old Password Passed, next up updating password");
        }
    };
    xmlhttp.open("POST", api_call, true);
    xmlhttp.setRequestHeader("Content-Type", "application/json");

    const data = [];
    for (let i = 0; i < old_pwd_txt.value.length; i++) {
        data.push(old_pwd_txt.value.charCodeAt(i));
    }

    xmlhttp.send(JSON.stringify({ password: data }));
    console.log("one");
    // update new password
    const xmlhttp2 = new XMLHttpRequest(); // new HttpRequest instance
    console.log("two");
    xmlhttp2.onreadystatechange = function() {
        check();
        console.log("STUFSSSS");
        if (xmlhttp2.readyState !== 4) {
            console.log("not ready");
            return;
        }
        if (xmlhttp2.status !== 200) {
            console.log("Failed to update Password");
            console.log(xmlhttp2.status);
            console.log(xmlhttp2.responseText);
            Command: toastr["error"]("Updating password failed.")

            toastr.options = {
                "closeButton": false,
                "debug": false,
                "newestOnTop": true,
                "progressBar": false,
                "positionClass": "toast-top-full-width",
                "preventDuplicates": true,
                "onclick": null,
                "showDuration": "300",
                "hideDuration": "1000",
                "timeOut": "5000",
                "extendedTimeOut": "1000",
                "showEasing": "swing",
                "hideEasing": "linear",
                "showMethod": "fadeIn",
                "hideMethod": "fadeOut"
            }
            var elem = document.getElementById('progress_div');
            elem.parentNode.removeChild(elem);

            var l = document.getElementById('alert_box');
            l.innerHTML = '<strong>Warning!</strong> Password update failed, could not update password. Please try again.'
            l.classList.remove("d-none");

            old_pwd_lbl.classList.remove("text-danger");
            old_pwd_lbl.classList.remove("text-success");
            old_pwd_txt.classList.remove("is-invalid");
            old_pwd_txt.classList.remove("is-valid");


            new_pwd_lbl.classList.remove("text-danger");
            new_pwd_lbl.classList.remove("text-success");
            new_pwd_txt.classList.remove("is-invalid");
            new_pwd_txt.classList.remove("is-valid");

            conf_pwd_lbl.classList.remove("text-danger");
            conf_pwd_lbl.classList.remove("text-success");
            conf_pwd_txt.classList.remove("is-invalid");
            conf_pwd_txt.classList.remove("is-valid");

            old_pwd_txt.disabled = false;
            new_pwd_txt.disabled = false;
            conf_pwd_txt.disabled = false;
            submit_btn.disabled = false;
            return;
        }
        else {
            Command: toastr["success"]("Password Updated successfully ")

            toastr.options = {
                "closeButton": true,
                "debug": false,
                "newestOnTop": true,
                "progressBar": false,
                "positionClass": "toast-top-full-width",
                "preventDuplicates": false,
                "onclick": null,
                "showDuration": "300",
                "hideDuration": "1000",
                "timeOut": "5000",
                "extendedTimeOut": "1000",
                "showEasing": "swing",
                "hideEasing": "linear",
                "showMethod": "fadeIn",
                "hideMethod": "fadeOut"
            }


            old_pwd_txt.disabled = false;
            old_pwd_txt.innerHTML = "";
            new_pwd_txt.disabled = false;
            new_pwd_txt.innerHTML = "";
            conf_pwd_txt.disabled = false;
            conf_pwd_txt.innerHTML = "";
            submit_btn.disabled = false;
            var l = document.getElementById('alert_box');
            l.classList.add("d-none");
            var p = document.getElementById('progress')
            p.style.width = '100%';
            console.log("New Password set");
            alert("Password updated");
        }
    };
    console.log("there");
    xmlhttp2.open("PUT", api_call, true);
    console.log("four");
    xmlhttp2.setRequestHeader("Content-Type", "application/json");
    console.log("five");
    var data2 = [];
    console.log("six");
    for (let i = 0; i < new_pwd_txt.value.length; i++) {
        data2.push(new_pwd_txt.value.charCodeAt(i));
    }
    console.log(data2);

    xmlhttp2.send(JSON.stringify({ password: data2 }));
    console.log("eight");
}
