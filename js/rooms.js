var buttons = document.querySelectorAll(".room_id_button")

for (var i = 0; i < buttons.length; i++) {
    buttons[i].addEventListener("click", (event) => {
        const url = "/room/" + event.target.value;
        window.open(url, '_blank');
    })
}

buttons = document.querySelectorAll(".delete_room_button")

for (var i = 0; i < buttons.length; i++) {
    buttons[i].addEventListener("click", (event) => {
        fetch("/rooms", {
            method: "POST",
            headers: {
                "contentType": "application/json",
                "accept": "application/json",
            },
            body: JSON.stringify({
                "room_id": event.target.value,
            })
        })
        .then( resp => {
            if (resp.ok) {
                resp.text().then( dat => {
                    const id = "row_" + event.target.value;
                    var rows = document.getElementsByName(id);
                    for (var i = 0; i < rows.length; i++) {
                        rows[i].remove();
                    }
                })
            } else {
                window.alert("Deleting romm #" + event.target.value + 
                    " failed!");
            }
        })
    })
}



document.getElementById("b_new").addEventListener("click", () => {
    window.open("/", "_self")
})

const date_cols = document.querySelectorAll(".date_col")
const comparison = '2000-01-01';

for (var i = 0; i < date_cols.length; i++) {
    const date = date_cols[i].textContent;
    if (date < comparison) {
        date_cols[i].textContent = ""
    }
}