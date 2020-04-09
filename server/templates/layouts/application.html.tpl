<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8">
    <title>{% block title %}Fascia{% endblock %}</title>
    <link rel="icon" href={{ "/lp/images/favicon.ico" | suffixAssetsUpdate }} type="image/vnd.microsoft.icon">
    <link href="https://cdnjs.cloudflare.com/ajax/libs/octicons/4.4.0/font/octicons.css" rel="stylesheet" type="text/css">
  </head>
  <body>
    <div id="app">
      {% block content %}{% endblock %}
    </div>
    <script src="https://use.fontawesome.com/080be9d465.js"></script>
    <script type="text/javascript" src={{ "js/main.js" | digestedAssets }}></script>
  </body>
</html>
