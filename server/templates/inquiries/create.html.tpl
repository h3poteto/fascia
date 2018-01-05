{% extends "layouts/not_login.html.tpl" %}

{% block title %}{{ title }}{% endblock %}

{% block content %}
<div class="inquiry">
  {% include "layouts/_login_header.html.tpl" %}
  <div class="main">
    <article>
      <div class="title">
        <h3>Contact</h3>
      </div>
      <div class="content">
        <h4>Thank you for your opinion. Please wait for a reply from us.</h4>
      </div>
    </article>
  </div>
</div>
{% endblock %}
