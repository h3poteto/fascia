{% extends "layout.html.tpl" %}

{% block title %}{{ title }}{% endblock %}

{% block content %}
<form action="/sign_up" method="post" role="form" name="sign_up" id="sign_up">
    <div class="form-group">
        <label for="email">Email:</label>
        <input class="form-control" name="email" type="email" />
    </div>

    <div class="form-group">
        <label for="password">Password:</label>
        <input class="form-control" name="password" type="password" />
    </div>

    <div class="form-group">
        <label for="password-confirm">Password:</label>
        <input class="form-control" name="password-confirm" type="password" />
    </div>

    <div class="form-action">
        <button class="btn" type="submit">SignUp</button>
    </div>
</form>
<p><a href="sign_in">SignIp</a></p>
{% endblock %}
