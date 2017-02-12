<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8">
    <title>{% block title %}Fascia{% endblock %}</title>
    <link href='https://fonts.googleapis.com/css?family=Raleway' rel='stylesheet' type='text/css'>
    <link href="https://cdnjs.cloudflare.com/ajax/libs/octicons/4.4.0/font/octicons.css" rel="stylesheet" type="text/css">
    <link rel="stylesheet" href={{ "/stylesheets/application-webview.css" | suffixAssetsUpdate }} media="all">

  </head>
  <body>
    <div id="content">
      {% block content %}{% endblock %}
    </div>
    <script src="https://use.fontawesome.com/080be9d465.js"></script>
    <script>
     (function(i,s,o,g,r,a,m){i['GoogleAnalyticsObject']=r;i[r]=i[r]||function(){
                             (i[r].q=i[r].q||[]).push(arguments)},i[r].l=1*new Date();a=s.createElement(o),
                             m=s.getElementsByTagName(o)[0];a.async=1;a.src=g;m.parentNode.insertBefore(a,m)
                             })(window,document,'script','https://www.google-analytics.com/analytics.js','ga');

     ga('create', 'UA-48724286-4', 'auto');
     ga('send', 'pageview');

    </script>
  </body>
</html>
