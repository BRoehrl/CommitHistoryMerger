$(".repoSelection").on('change', function(event) {
  $.post('./repositories/'+this.name+"/"+this.value)
});
