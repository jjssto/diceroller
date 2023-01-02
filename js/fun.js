export { 
    hasFocus, 
    highlightOwnRolls,
    addCol,
    init,
    initReset
};

var hasFocus;
var highlightOwnRolls;

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

function insertOption(sel, text, value) {
    var element = document.createElement("option");
    element.value = value
    element.textContent = text
    sel.appendChild(element);
}


function insertColorOptions() {
    
    // get setting form, if not present => exit
    var sel = document.getElementById("f_setting_color");
    if (sel == null) return

    // create default option
    insertOption(sel, "-", "-")
    
    // get and set the additional color options
    fetch("/res/colors.json", {
        method: "GET"
    })
    .then( resp => {
        resp.json().then(
            data => {
               for (var i in data) {

                    insertOption(sel, data[i].text, data[i].code)
               }
            }
        )
    })
}



function roomSettingForm() {
    //insertColorOptions()
    var settingsForm = document.getElementById("f_setting");
    settingsForm.addEventListener("submit", function (event) {
        event.preventDefault();
        var roomName = document.getElementById("f_setting_name").value;
        var colorSelect = document.getElementById("f_setting_color");
        var color = colorSelect.options[colorSelect.selectedIndex].value;
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
    document.querySelector(".header").style.backgroundColor = color;
    document.querySelector(".footer").style.backgroundColor = color;
    document.querySelectorAll("thead").forEach(element => {
       element.style.backgroundColor = color; 
    });
    document.querySelectorAll("th").forEach(element => {
       element.style.backgroundColor = color; 
    });

    document.getElementById("f_setting_color").value = color;
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

function setLink() {
    var link = document.getElementById("a_link");
    const loc = window.location.href;
    link.href = loc
    link.textContent = loc
}

function addCol(row, text) {
    const col = row.insertCell()
    col.appendChild(document.createTextNode(text))
}



function insertRolls(createRow, data_raw) {
    const data = JSON.parse(data_raw)
    const tbody = document.getElementById("tbody_rolls")
    const first_row = tbody.firstChild
    var roll_id
    if (first_row) {
        roll_id = first_row.firstChild.firstChild.textContent
    } else {
        roll_id = -1
    }
    var counter = 0
    var lastRow
    for (let i in data) {
        const drow = data[i]
        sessionStorage.setItem("last_roll", i)
        if (i > roll_id) { 
           const row = createRow(drow, i)
           lastRow = tbody.appendChild(row)
        }
        counter++
    }
    if (lastRow) {
        lastRow.scrollIntoView(true)
    }
}

function getRolls(createRow) {
    if (!hasFocus) {
        return
    }
    var target = location.href.replace("room/", "rolls/");
    const last_roll = sessionStorage.getItem("last_roll");
    var offsetStr = sessionStorage.getItem("ts_offset");
    if (offsetStr == null || offsetStr.length == 0) {
        const date = new Date();
        const offset = -1 * date.getTimezoneOffset() * 60; 
        offsetStr = offset.toString();
        sessionStorage.setItem("ts_offset", offsetStr)
    }
    if (last_roll != "") {
        target += "/" + last_roll
    }
    fetch(target, {
        method: "GET",
        headers: {
            "ts_offset": offsetStr
        },
    })
    .then(response => {
        response.text().then(data => {
                insertRolls(createRow, data);
        })
    });
}

function setHighlightOwnRolls() {
    const element = document.getElementById("i_highlight_own_rolls")
    if (element == null) {
        return
    }
    if (element.checked) {
        highlightOwnRolls = true;
    } else {
        highlightOwnRolls = false;
    }
}

function init(createRow) {

    window.addEventListener("DOMContentLoaded", () => {
        setHighlightOwnRolls()
        setColor()
        detectFocus()
        settingVisibility()
        roomSettingForm()
        setLink()

    })
    window.addEventListener("load", () => {
        sessionStorage.setItem("last_roll", "");
        getRolls(createRow);
    });

    window.setInterval(() => {getRolls(createRow)}, 1000);
}

function initReset(reset) {
    document.getElementById("b_reset").addEventListener("click", () => {
        reset();
    });
}