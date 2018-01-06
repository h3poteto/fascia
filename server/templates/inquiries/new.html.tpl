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
      <div class="flash flash-error">{{ error }}</div>
      <div class="content">
        <div class="contact-board">
          <form action="/inquiries" method="post" role="form" name="inquiries" id="inquirires" class="pure-form pure-form-stacked">
            <fieldset>
              <div class="pure-control-group control-group">
                <label for="message">Message</label>
                <textarea class="form-control" name="message"></textarea>
              </div>
              <div class="pure-control-group control-group">
                <label for="email">Email</label>
                <input class="form-control" name="email" type="email" placeholder="Your email address" />
              </div>
              <div class="pure-control-group control-group">
                <label for="name">Name</label>
                <input class="form-control" name="name" placeholder="Your name" />
              </div>
              <div class="pure-controls control-group contact-control">
                <button class="pure-button pure-button-primary contact-button" type="submit">Submit</button>
              </div>
            </fieldset>
          </form>
        </div>
      </div>
    </article>
  </div>
</div>
{% endblock %}
