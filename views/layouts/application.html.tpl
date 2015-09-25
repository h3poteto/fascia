<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8">
    <title>{% block title %}Fascia{% endblock %}</title>
    <link rel="stylesheet" href="/stylesheets/pure-min.css" media="all">
    <link rel="stylesheet" href="/stylesheets/octicons.css" media="all">
    <link rel="stylesheet" href="/stylesheets/application.css" media="all">

  </head>
  <body>
    <div id="top">
    </div>
    <div id="content">
      {% block content %}{% endblock %}
    </div>
    <script type="text/javascript" src="/javascripts/bundle.js"></script>
  </body>
</html>
