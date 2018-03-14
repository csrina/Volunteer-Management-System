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



    if (new_pwd_txt.value == conf_pwd_txt.value) {
        conf_message.classList.add("d-none");
        conf_pwd_lbl.classList.remove("text-danger");
        conf_pwd_lbl.classList.add("text-success");
        conf_pwd_txt.classList.remove("is-invalid");
        conf_pwd_txt.classList.add("is-valid");

    }
    else if (new_pwd_txt.value != conf_pwd_txt.value && conf_pwd_txt.value.length > 0) {
        conf_message.classList.remove("d-none");
        conf_pwd_lbl.classList.add("text-danger");
        conf_pwd_lbl.classList.remove("text-success");
        conf_pwd_txt.classList.add("is-invalid");
        conf_pwd_txt.classList.remove("is-valid");
    }
    if (conf_pwd_txt.value.length == 0) {
        conf_message.classList.add("d-none");
        conf_pwd_lbl.classList.remove("text-danger");
        conf_pwd_lbl.classList.remove("text-success");
        conf_pwd_txt.classList.remove("is-invalid");
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
    var l = getElementById('alert');
    l.parentNode.removeChild(l);

    var ready = check();
    if (ready == false) {
        return 0;
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
        "positionClass": "toast-top-right",
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

    var old_pwd_txt = document.getElementById('new_pwd');
    var old_pwd_lbl = document.getElementById('new_pwd_lbl');

    var new_pwd_txt = document.getElementById('new_pwd');
    var new_pwd_lbl = document.getElementById('new_pwd_lbl');


    var conf_pwd_txt = document.getElementById('confirm_pwd');
    var conf_pwd_lbl = document.getElementById('confirm_pwd_lbl');


    // check old password
    let api_call = "/api/v1/passwords"

    const xmlhttp = new XMLHttpRequest();
    xmlhttp.onreadystatechange() = function() {
        if (this.readyState !== 4) return;
        if (this.status !== 200) {
            console.log(this.status);
            Command: toastr["error"]("Updating password failed.")

            toastr.options = {
                "closeButton": false,
                "debug": false,
                "newestOnTop": true,
                "progressBar": false,
                "positionClass": "toast-top-right",
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
            return elem.parentNode.removeChild(elem);
            location = document.getElementById('submit_div');
            localStorage.insertAdjacentHTML('beforebegin', '<div id="alert" class="alert alert-warning"> <strong>Warning!</strong> Password update failed, could not authenticate user.</div>')
        }
        else {
            var p = getElementById('progress')

        }
    }



    // update new password


}
