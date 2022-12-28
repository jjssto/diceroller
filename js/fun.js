export { detectFocus, hasFocus, settingVisibility, roomSettingForm, setColor };

var hasFocus;

function detectFocus() {
    hasFocus = true
    window.addEventListener('focus', function (event) {
        hasFocus = true
        window.addEventListener('blur', function (event) {
            hasFocus = false
        })
    })
}

function settingVisibility() {
    var element = document.getElementById("f_setting");
    if (sessionStorage.getItem("is_owner") == "true") {
        element.style.visibility = 'visible';
    } else {
        element.style.visibility = 'hidden';
    }
}

function roomSettingForm() {
    var settingsForm = document.getElementById("f_setting");
    settingsForm.addEventListener("submit", function (event) {
        event.preventDefault();
        var roomName = document.getElementById("f_setting_name").value;
        var color = document.getElementById("f_setting_color").value;
        var roomId = sessionStorage.getItem("room_id");
        fetch("/roomSettings", {
            method: "POST",
            headers: {
                "Accept": "application/json",
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                "room_id": roomId,
                "room_name": roomName,
                "color": color
            })
        })
        .then(res => { 
            window.location.reload()
        })
    })
}

function changeColor(color) {
    //document.getElementById("title_line").style.background = color;
    document.getElementById("title_tag").style.background = color;
    document.getElementById("table_head").style.background = color;
}

function setColor() {
    window.addEventListener("load", function (event) {
        var color = this.sessionStorage.getItem("color");
        if (color.length > 0) {
            changeColor(color);
        }
    })
    window.addEventListener("reset", function (event) {
        var color = this.sessionStorage.getItem("color");
        if (color.length > 0) {
            changeColor(color);
        }
    })

}