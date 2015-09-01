{% extends "layout.html.tpl" %}

{% block title %}{{ title }}{% endblock %}

{% block content %}
<form action="/sign_in" method="post" role="form" name="sign_in" id="sign_in">
    <div class="form-group">
        <label for="email">Email:</label>
        <input class="form-control" name="email" type="email" />
    </div>

    <div class="form-group">
        <label for="password">Password:</label>
        <input class="form-control" name="password" type="password" />
    </div>

    <div class="form-action">
        <button class="btn" type="submit">SignIn</button>
    </div>
</form>
<p><a href="sign_up">SignUp</a></p>
{% endblock %}
