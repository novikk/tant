<html>
    <head>
        <link rel="stylesheet" href="static/css/style.css">
        <link href="https://fonts.googleapis.com/css?family=Heebo" rel="stylesheet">
    </head>
    <body>
        <div class="form">
            <div class="top_banner">
                <img src="static/img/Formulario/Header2.jpeg" alt="">
            </div>

            <form action="/set_info" method="POST">
                Gender <br>
                <div class="outer">
                    <div class="inner"><input type="radio" name="gender" value="male" checked> &nbsp;&nbsp;Male</div>
                    <div class="inner" style="margin-left: 10%"><input type="radio" name="gender" value="female"> &nbsp;&nbsp;Female</div>
                </div>


                Age <br>
                <input name="age" type="text"><br><br>

                Country <br>
                <input name="country" type="text"><br><br>

                <div class="form-btn">
                    <input type="image" src="static/img/Formulario/Button.png" alt="">
                </div>
            </form>
        </div>
    </body>
</html>