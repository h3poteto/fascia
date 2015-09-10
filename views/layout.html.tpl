<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>{% block title %}Fascia{% endblock %}</title>
    <link rel="stylesheet" href="/stylesheets/pure-min.css" media="all">
</head>
<body>
    <div id="content">
        {% block content %}{% endblock %}
    </div>
    <script src="/javascripts/bundle.js">
</body>
