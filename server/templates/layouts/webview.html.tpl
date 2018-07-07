<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>{% block title %}Fascia{% endblock %}</title>
    <link rel="icon" href={{ "/images/favicon.ico" | suffixAssetsUpdate }} type="image/vnd.microsoft.icon">
    <link href='https://fonts.googleapis.com/css?family=Raleway' rel='stylesheet' type='text/css'>
    <link href="https://cdnjs.cloudflare.com/ajax/libs/octicons/4.4.0/font/octicons.css" rel="stylesheet" type="text/css">
    <link rel="stylesheet" href={{ "stylesheets/application-webview.css" | digestedAssets }} media="all">

  </head>
  <body>
    <div id="content">
      {% block content %}{% endblock %}
    </div>
    <script src="https://use.fontawesome.com/080be9d465.js"></script>
  </body>
</html>
