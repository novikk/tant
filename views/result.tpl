<html>
    <head>
        <link rel="stylesheet" href="static/css/style.css">
        <link href="https://fonts.googleapis.com/css?family=Heebo" rel="stylesheet">
    </head>
    <body>
        <div class="results">
            <div class="top_banner">
                <img src="static/img/Individuales/Header.png" alt="">
            </div>

            <div class="grid">
                <div class="grid_item_left">
                    <a href="{{.art1.ShopURL}}"><img src="{{(index .art1.Media.Images 0).LargeURL}}" alt="">
                    {{.art1.Brand.Name}}</a>
                </div>

                <div class="grid_item_right">
                    <a href="{{.art2.ShopURL}}"><img src="{{(index .art2.Media.Images 0).LargeURL}}" alt="">
                    {{.art2.Brand.Name}}</a>
                </div>

                <div class="grid_item_left" style="margin-top: 2%">
                    <a href="{{.art3.ShopURL}}"><img src="{{(index .art3.Media.Images 0).LargeURL}}" alt="">
                    {{.art3.Brand.Name}}</a>
                </div>

                <div class="grid_item_right" style="margin-top: 2%">
                    <a href="{{.art4.ShopURL}}"><img src="{{(index .art4.Media.Images 0).LargeURL}}" alt="">
                    {{.art4.Brand.Name}}</a>
                </div>
            </div>

            <div class="recom_btn">
                <img src="static/img/Individuales/Button.png" alt="" onclick="javascript:location.reload()">
            </div>
        </div>
    </body>
</html>