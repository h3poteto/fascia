<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8">
    <title>{% block title %}Fascia{% endblock %}</title>
    <link href='https://fonts.googleapis.com/css?family=Raleway' rel='stylesheet' type='text/css'>
    <link rel="stylesheet" href={{ "/stylesheets/pure-min.css" | suffixAssetsUpdate }} media="all">
    <link rel="stylesheet" href={{ "/stylesheets/octicons.css" | suffixAssetsUpdate }} media="all">
    <link rel="stylesheet" href={{ "/stylesheets/font-awesome.css" | suffixAssetsUpdate }} media="all">
    <link rel="stylesheet" href={{ "/stylesheets/application-webview.css" | suffixAssetsUpdate }} media="all">

  </head>
  <body>
    <div id="content">
      {% block content %}{% endblock %}
    </div>
  </body>
</html>
