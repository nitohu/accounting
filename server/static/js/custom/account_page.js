$(document).ready(function() {
    let btns = document.getElementsByClassName("deleteEntry")

    for (let i = 0; i < btns.length; i++) {
        let btn = btns[i]
        btn.addEventListener("click", deleteAccount)
    }

    function deleteAccount(e) {
        let id = $(e.target).attr("account-id")
        data = {"ID": Number.parseInt(id)}
        
        let xhr = new XMLHttpRequest()
        xhr.open("DELETE", "/api/accounts/delete", true)
        xhr.setRequestHeader("Content-Type", "application/json")
        xhr.onreadystatechange = function() {
            if (this.readyState == 4) {
                let m = this.response.replace(/'/g, '"')
                let msg = JSON.parse(m)
                if (this.status == 400) {
                    console.error(msg.error)
                } else if (this.status == 200) {
                    window.location.reload()
                    console.log(msg.success)
                } else if (this.status == 403) {
                    console.warn(msg.error)
                } else {
                    console.error(this.response)
                }
            }
        }
        xhr.send(JSON.stringify(data))
    }
})
