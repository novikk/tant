<html>
    <head>
        <link rel="stylesheet" href="static/css/style.css">
        <link href="https://fonts.googleapis.com/css?family=Heebo" rel="stylesheet">
    </head>
    <body>
        <div class="logins">
            <div class="top_banner">
                <img src="static/img/Logins/Header.png" alt="">
            </div>

            <div class="logins_msg">
                Log in to your favourite <br> accounts
            </div>

            {{ if eq .deezerUser.Email "" }}
            <div class="logins_btn1">
                <a id="deezerButton" onclick="block_none()" href="/login/deezer" ><img src="static/img/Logins/Deezer.png" alt=""></a>
                <a id="deezerGif" style="font-size:30px" class="hide"><img src="static/img/ajax-loader.gif" alt=""><br>Logging into your deezer account</a>

            </div>
            {{ else }}
                <div class="logins_btn2">
                    Logged in Deezer as {{.deezerUser.Email}}
                </div>
            {{ end }}

            {{ if eq .instagramUser.Name "" }}
                <div class="logins_btn2">
                    <a href="/login/instagram" id=""><img src="static/img/Logins/Instagram.png" alt=""></a>
                </div>
            {{ else }}
                <div class="logins_btn2">
                    Logged in Instagram as {{.instagramUser.Name}}
                </div>
            {{ end }}

            <div class="logins_magic">
                <a href="/recommendations"><img src="static/img/Logins/Button.png" alt=""></a>
            </div>
        </div>
    </body>

    <script> 
        function block_none(){
            document.getElementById('deezerButton').classList.add('hide');
            document.getElementById('deezerGif').classList.add('show');
        }
    </script>
</html>