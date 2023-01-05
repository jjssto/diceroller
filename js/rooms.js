const buttons = document.querySelectorAll(".room_id_button")

for (var i = 0; i < buttons.length; i++) {
    buttons[i].addEventListener("click", (event) => {
        const url = "/room/" + event.target.value;
        window.open(url, '_blank');
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