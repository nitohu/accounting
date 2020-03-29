$(document).ready(function() {

    function store(value) {
        let expiryDate = new Date(Date.now())
        expiryDate.setDate(expiryDate.getDate() + 30)

        let formattedValue = ""
        for (let i = 0; i < value.length; i++) {
            formattedValue += value[i]
            if(i < value.length - 1) {
                formattedValue += ","
            }
        }
        
        let cookieStr = "theme="+formattedValue+"; expires="+ expiryDate.toGMTString() +"; path=/;"
        document.cookie = cookieStr
    }

    function read() {
        let c = document.cookie
        console.log(c)
        let cookies = c.split(" ")

        for(let i = 0; i < cookies.length; i++) {
            let [k, v] = cookies[i].split("=")

            if (k != "theme") {
                continue
            }

            if (v[v.length - 1] == ";") {
                v = v.slice(0, -1)
            }

            return v.split(",")
        }
    }

    function getCurrentSettings() {
        console.log("getCurrentSettings")
        let c = $("body").prop("class")
        
        let classes = c.split(" ")
        console.log(classes)

        let themeSettings = []
        for (let i = 0; i < classes.length; i++) {
            let setting = ""
            if (classes[i].includes("-")) {
                setting = classes[i].split("-")[1]
            } else {
                setting = classes[i]
            }
            themeSettings.push(setting)
        }
        console.log(themeSettings)
        store(themeSettings)
        updateUI(themeSettings)
    }

    function getColorSettings() {
        // Wait 100ms before executing the function
        // Otherwise it can happen that the old color setting gets stored into the cookie
        setTimeout(getCurrentSettings, 100)
    }

    function updateUI(data) {
        if (data == undefined) {
            return
        }
        console.log(data.length)
        if(data.length > 0) {
            $("body").removeClass()
        }

        for(let i = 0; i < data.length; i++) {
            let item = data[i]

            if (item == "toggle") {
                item = "ls-toggle-menu"
                $("#checkbox2").prop("active")
            } else {
                item = "theme-"+item
            }

            $("body").addClass(item)
        }
    }

    console.log(read())
    // getCurrentSettings()
    updateUI(read())

    $(".themeSetting").change(getCurrentSettings)
    $(".themeButton").click(getCurrentSettings)
    $(".themeColorBtn").click(getColorSettings)
})