/* Wait for DOM to load, then get the gauges */
document.addEventListener("DOMContentLoaded", () => {
    // Make request for the donation page data
    $.ajax({
        url: "/api/v1/donate",
        type: 'GET',
        contentType:'json',
        success: (data => createWidgets(data)),
        error: (xhr => makeToast("error", "Failed to retrieve page data: " + xhr.responseText)),
        dataType: 'json'
    });
})

function sendDonation() {
    if (!validateForm()) {
        makeToast("warning", "Donation amount cannot exceed hours available!");
    }

    let doneeIndex = document.getElementById("DoneeSelect").selectedIndex;
    let data = {
        amount: $("#AmountField").attr("value"),
        donee: document.getElementById("DoneeSelect").options[doneeIndex].value
    };
    $.ajax({
        url: "/api/v1/dashboard",
        type: 'GET',
        contentType:'json',
        data: JSON.stringify(data),
        success: function(data) {
            makeToast("success", "Donation sent! Donation ID: " + data.id);
            refreshWidgets(data);
        },
        error: (xhr => makeToast("error", "Failed to send donation: " + xhr.responseText)),
        dataType: 'json'
    });
}

function createWidgets(data) {
    data.families.forEach(fam => {
        $("#DoneeSelect").append("<option value='''" + fam.familyID + "' >" + fam.familyName + "</option>");
    });
    document.getElementById("hoursAvail-text").innerHTML = (data.hoursAvail > 0) ? data.hoursAvail : 0;
}

// Changes appearance and returns a bool (isValid? true if valid)
function validateForm() {
    if ($("#AmountField").html() > $("#hoursAvail-text").html()) {
        $("#AmountValidFB").css("visibility", "hidden");
        $("#AmountInvalidFB").css("visibility", "visible");
        $("#AmountField").removeClass("is-valid").addClass("is-invalid");
        $("#amountFG").removeClass("has-success").addClass("has-danger");
        return false;
    }
    $("#AmountInvalidFB").css("visibility", "hidden");
    $("#AmountValidFB").css("visibility", "visible");
    $("#AmountField").removeClass("is-invalid").addClass("is-valid");
    $("#amountFG").removeClass("has-danger").addClass("has-success");
    return true;
}
