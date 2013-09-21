function focusSelect(field) {
  field.focus()
  field.select()
}

function handlerAccount() {
  focusSelect(document.form.email)
}

function handlerAdd() {
  focusSelect($('table.catalog tr td.keyword input.key').first())

  $('table.catalog tr td.icon, table.catalog tr td.name').click(function() {
    focusSelect($(this).parent().find('input.key'))
  })
}

function handlerHome() {
  focusSelect(document.form.q)
}

function handlerDelete() {
  var form = document.form

  if (form.key.type == 'hidden') {
    focusSelect(form.delete)
  } else {
    focusSelect(form.key)
  }
}

function handlerEdit() {
  var form = document.form

  if (form['newkey'].value == '') {
    focusSelect(form['newkey'])
  } else {
    focusSelect(form.body)
  }
}

function handlerHelpGlossary() {
  var highlight = function(value) {
    if (!value) return
    $('dl.glossary dt a').css('background-color', '')
    $('dl.glossary dt a[name=' + value.substring(1) + ']').attr('style', 'background-color:yellow')
  }
  highlight(window.location.hash)

  if ('onhashchange' in window) {
      window.onhashchange = function() {
          highlight(window.location.hash)
      }
  } else {
      var storedHash = window.location.hash
      window.setInterval(function () {
          if (window.location.hash != storedHash) {
              storedHash = window.location.hash
              highlight(storedHash)
          }
      }, 100)
  }
}

function handlerLogin() {
  focusSelect(document.form.email)
}

function handlerList() {
  var e = $('td.key a').first()
  if (e) focusSelect(e)
}

function handlerRecover() {
  var form = document.form

  if (form.token) {
    focusSelect(form['newpassword'])
  } else {
    focusSelect(form.email)
  }
}

function handlerRename() {
  var form = document.form

  if (form['srckey'].type == 'text') {
    focusSelect(form['srckey'])
  } else if (form['dstkey'].type == 'text') {
    focusSelect(form['dstkey'])
  } else {
    focusSelect(form.rename)
  }
}

function handlerSignup() {
  focusSelect(document.form.email)
}

$(document).ready(function() {
  switch (location.pathname) {
    case '/':
      handlerHome()
      break
    case '/run':
      handlerHome()
      break
    case '/account':
      handlerAccount()
      break
    case '/add':
      handlerAdd()
      break
    case '/delete':
      handlerDelete()
      break
    case '/edit':
      handlerEdit()
      break
    case '/help/glossary':
      handlerHelpGlossary()
      break
    case '/list':
      handlerList()
      break
    case '/login':
      handlerLogin()
      break
    case '/recover':
      handlerRecover()
      break
    case '/rename':
      handlerRename()
      break
    case '/signup':
      handlerSignup()
      break
  }

  // Ctrl+enter submits form
  $('form').keydown(function (e) {
    if (e.ctrlKey && e.keyCode == 13) { $(e.target).closest('form').submit() }
  })

  // Esc navigates to homepage
  $(document).keyup(function(e) {
    if (e.keyCode == 27) { window.location = '/' }
  })
})
