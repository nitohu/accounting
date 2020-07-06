$(document).ready(function() {
    const rs = $("#rightsidebar")
    const b = $("body")
    const colors = ["purple", "blue", "cyan", "green", "orange", "blush"]

    updateUI(_getThemeCookie())

    
    $(".page-loader-wrapper").fadeOut()
    
    // Left Sidebar: Menuitems with dropdown
    let menuitems = $(".menu-toggle")
    for(let i = 0; i < menuitems.length; i++) {
        let e = $(menuitems[i])
        let m = e.parent().find(".ml-menu")
        let l = m.find("li")

        for(let i = 0; i < l.length; i++) {
            if (l[i].className.includes("active")) {
                m.fadeToggle()
            }
        }

        e.click(function() {
            m.fadeToggle()
            e.toggleClass("toggled")
        })
    }

    // Toggles left sidebar
    $(".ls-toggle-btn").click(function() {
        b.toggleClass("ls-toggle-menu")
        saveCurrentTheme()
    })

    // Toggles right sidebar
    $(".right_icon_toggle_btn").click(function() {
        b.toggleClass("right_icon_toggle")
        saveCurrentTheme()
    })

    // Settings panel
    $("#settingsBtn").click(function() {
        rs.toggleClass("open")
        rs.hasClass("open") ? $(".overlay").fadeIn() : $(".overlay").fadeOut()
    })

    // Remove settings panel when clicking somewhere
    $(document).click(function(e) {
        let t = $(e.target)
        if ( rs.hasClass("open") && (!t.hasClass("js-right-sidebar")
            && !t.parent().hasClass("js-right-sidebar") && rs.find(t).length == 0)) {
            $(".overlay").fadeOut();
            rs.removeClass("open");
        }
    })

    // Settings: Dark/Bright Theme
    $(".themeSetting").click(function() {
        if ($("#darktheme").is(":checked")) {
            b.addClass("theme-dark")
        } else {
            b.removeClass("theme-dark")
        }
        saveCurrentTheme()
    })

    // Settings: Theme Color
    $(".themeColorBtn").click(function(e) {
        let p = $(e.target).parent()
        let data = getCurrentThemeData()
        data = data.filter((v) => !colors.includes(v))
        
        let c = p.attr("data-theme")
        data.push(c)

        console.log("Data: ")
        console.log(data)

        saveTheme(data)
    })

    // Settings: RTL button
    $("#checkbox1").click(function() {
        $("body").toggleClass("rtl")

        $("#checkbox1").prop("checked", !$("#checkbox1").prop("checked"))

        saveCurrentTheme()
    })

    // Settings: Mini Sidebar button
    $("#checkbox2").click(function() {
        $("body").toggleClass("ls-toggle-menu")

        $("#checkbox2").prop("checked", !$("#checkbox2").prop("checked"))

        saveCurrentTheme()
    })
})