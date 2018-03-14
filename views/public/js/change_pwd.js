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

function setAlertText(txt) {
    var alrt = document.getElementById('alert_box');
    alrt.innerHTML = txt;
}

function updateStatus(disable, neutral, isFailedPassword, isFailedReq) {
    // status == 1 = disable fields
    // status == 0 = enable fields
    var old_pwd_txt = document.getElementById('old_pwd');
    var old_pwd_lbl = document.getElementById('old_pwd_lbl');
    var old_message = document.getElementById('old_message');

    var new_pwd_txt = document.getElementById('new_pwd');
    var new_pwd_lbl = document.getElementById('new_pwd_lbl');
    var new_message = document.getElementById('new_message');

    var conf_pwd_txt = document.getElementById('confirm_pwd');
    var conf_pwd_lbl = document.getElementById('confirm_pwd_lbl');
    var conf_message = document.getElementById('conf_message');

    var l = document.getElementById('alert_box');
    var submit_btn = document.getElementById('change_pwd_btn');

    if (disable) {
        old_pwd_txt.disabled = true;
        new_pwd_txt.disabled = true;
        conf_pwd_txt.disabled = true;
        submit_btn.disabled = true;
    }
    else {
        old_pwd_txt.disabled = false;
        new_pwd_txt.disabled = false;
        conf_pwd_txt.disabled = false;
        submit_btn.disabled = false;
    }

    if (neutral) {
        updateDisplay(old_pwd_lbl, old_pwd_txt, old_message, 0);
        updateDisplay(new_pwd_lbl, new_pwd_txt, new_message, 0);
        updateDisplay(conf_pwd_lbl, conf_pwd_txt, conf_message, 0);
    }

    if (isFailedPassword) {
        old_pwd_lbl.classList.add("text-danger");
        old_pwd_txt.classList.add("is-invalid");
    }

    if (isFailedReq) {
        l.classList.remove("d-none");
    }
    else {
        l.classList.add("d-none");
    }


}

function eraseTextFields() {
    var old_pwd_txt = document.getElementById('old_pwd');
    var new_pwd_txt = document.getElementById('new_pwd');
    var conf_pwd_txt = document.getElementById('confirm_pwd');
    old_pwd_txt.value = "";
    new_pwd_txt.value = "";
    conf_pwd_txt.value = "";
}

function showToaster(type, msg) {

    Command: toastr[type](msg);

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
    };
}

function addProgressBar() {
    var location = document.getElementsByClassName('container');
    var progress = location[0];
    progress.insertAdjacentHTML('afterbegin',
        '<div class="progress" id="progress_div" style="margin-bottom: 20px"> <div id="progress" class="progress-bar progress-bar-striped progress-bar-animated" role="progressbar" aria-valuenow="75" aria-valuemin="0" aria-valuemax="100" style="width: 25%"></div><br></div>')
}

function updateProgress(precent) {
    var p = document.getElementById('progress')
    p.style.width = precent;
}

function removeProgress() {
    var elem = document.getElementById('progress_div');
    elem.parentNode.removeChild(elem);
}

function updateDisplay(lbl, txt, msg, state) {
    var submit_btn = document.getElementById('change_pwd_btn');
    // state == 0 neutral
    // state == 1 pass/green
    // state == 2 fail/red
    if (state == 2) {
        msg.classList.remove("d-none");
        lbl.classList.add("text-danger");
        lbl.classList.remove("text-success");
        txt.classList.add("is-invalid");
        txt.classList.remove("is-valid");
        submit_btn.disabled = true;
    }
    if (state == 1) {
        msg.classList.add("d-none");
        lbl.classList.remove("text-danger");
        lbl.classList.add("text-success");
        txt.classList.remove("is-invalid");
        txt.classList.add("is-valid");
        submit_btn.disabled = false;
    }
    if (state == 0) {
        msg.classList.add("d-none");
        lbl.classList.remove("text-danger");
        lbl.classList.remove("text-success");
        txt.classList.remove("is-invalid");
        txt.classList.remove("is-valid");
        submit_btn.disabled = false;
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

    // state == 0 neutral
    // state == 1 pass/green
    // state == 2 fail/red

    // tests new password length
    if (new_pwd_txt.value.length < 8) {
        updateDisplay(new_pwd_lbl, new_pwd_txt, new_message, 2);
    }
    else {
        updateDisplay(new_pwd_lbl, new_pwd_txt, new_message, 1);

    }

    // --------------------------------- tests new password
    if (new_pwd_txt.value == conf_pwd_txt.value &&
        new_pwd_txt.value.length != 0 &&
        conf_pwd_txt.value.length != 0) {

        updateDisplay(conf_pwd_lbl, conf_pwd_txt, conf_message, 1);
    }
    else if (new_pwd_txt.value != conf_pwd_txt.value &&
        conf_pwd_txt.value.length > 0) {

        conf_message.innerText = "Needs to match new password";
        updateDisplay(conf_pwd_lbl, conf_pwd_txt, conf_message, 2);

    }
    else if (new_pwd_txt.value != conf_pwd_txt.value &&
        conf_pwd_txt.value.length == 0) {

        conf_message.innerText = "Cannot be empty and needs to match new password";
        updateDisplay(conf_pwd_lbl, conf_pwd_txt, conf_message, 2);

    }
    else if (new_pwd_txt.value.length == 0 &&
        conf_pwd_txt.value.length == 0) {

        conf_message.innerText = "Cannot be empty";
        updateDisplay(conf_pwd_lbl, conf_pwd_txt, conf_message, 2);

    }

    // --------------------------------- test old password
    if (old_pwd_txt.value.length > 0) {
        updateDisplay(old_pwd_lbl, old_pwd_txt, old_message, 0);

    }
    else {
        updateDisplay(old_pwd_lbl, old_pwd_txt, old_message, 2);
    }


    if (old_pwd_txt.value.length == 0 || conf_pwd_txt.value.length == 0 ||
        (new_pwd_txt.value != conf_pwd_txt.value && conf_pwd_txt.value.length > 0) ||
        new_pwd_txt.value.length < 8) {
        return false;
    }
    else {
        return true;
    }
}

function change_pwd() {
    //disable fields, set fields to neutral, isFailedPassword, is FailedReg 
    updateStatus(true, true, false, false);

    var ready = check();
    if (ready == false) {
        updateStatus(false, false, false, false);
        return;
    }
    addProgressBar();
    showToaster("info", "Updating Password ......");

    // check old password
    var old = document.getElementById('old_pwd').value;
    const api_call = "/api/v1/passwords";
    const xmlhttp = new XMLHttpRequest(); // new HttpRequest instance
    xmlhttp.onreadystatechange = function() {
        check();
        if (xmlhttp.readyState !== 4) { return; }
        if (xmlhttp.status !== 202) {
            console.log("Failed to auth");
            console.log(xmlhttp.status);
            console.log(xmlhttp.responseText);
            showToaster("error", "Updating password failed.");
            removeProgress();
            setAlertText('<strong>Warning!</strong> Password update failed, could not authenticate user. Please enter correct password.')
            updateStatus(false, false, true, true);
            return;
        }
        else {
            updateStatus(false, false, false, false);
            updateProgress('60%');
            console.log("User authenticated");
            alert("Old Password Passed, next up updating password");
        }
    };
    xmlhttp.open("POST", api_call, true);
    xmlhttp.setRequestHeader("Content-Type", "application/json");

    const data = [];
    for (let i = 0; i < old.length; i++) {
        data.push(old.charCodeAt(i));
    }

    xmlhttp.send(JSON.stringify({ password: data }));

    // update new password
    var newp = document.getElementById('new_pwd').value;
    const xmlhttp2 = new XMLHttpRequest(); // new HttpRequest instance
    xmlhttp2.onreadystatechange = function() {
        check();
        if (xmlhttp2.readyState !== 4) { return; }
        if (xmlhttp2.status !== 200) {
            console.log("Failed to update Password");
            console.log(xmlhttp2.status);
            console.log(xmlhttp2.responseText);
            showToaster("error", "Updating password failed.");
            removeProgress();
            setAlertText('<strong>Warning!</strong> Password update failed, could not update password. Please try again.')
            updateStatus(false, true, false, true);
            return;
        }
        else {
            showToaster("success", "Password Updated successfully ");
            updateStatus(false, true, false, false);
            updateProgress('100%');
            console.log("New Password set");
            alert("Password updated");
            eraseTextFields();
            removeProgress();
        }
    };
    xmlhttp2.open("PUT", api_call, true);
    xmlhttp2.setRequestHeader("Content-Type", "application/json");
    var data2 = [];
    for (let i = 0; i < newp.length; i++) {
        data2.push(newp.charCodeAt(i));
    }
    console.log(data2);

    xmlhttp2.send(JSON.stringify({ password: data2 }));
}
