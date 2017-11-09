package services

import (
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/mysql"
  "github.com/stanisdev/juicy-blog/models"
  "fmt"
)

func ConnectToDatabase(c *Config) *gorm.DB  {
  connectParams := c.DbUser + ":" + c.DbPass + "@/" + c.DbName + "?charset=utf8&parseTime=True&loc=Local"
  con, err := gorm.Open("mysql", connectParams)
  if err != nil {
    panic("failed to connect database")
  }
  con.LogMode(true)
  return con
}

func DatabaseMigrate(c *Config) {
  con := ConnectToDatabase(c)
  con.AutoMigrate(&models.User{}, &models.Article{}, &models.Subscriber{}, &models.NewArticlesSubscriber{})
  con.Model(&models.Article{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
  con.Model(&models.Subscriber{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
  con.Model(&models.Subscriber{}).AddForeignKey("subscriber_id", "users(id)", "CASCADE", "CASCADE")
  con.Model(&models.NewArticlesSubscriber{}).AddForeignKey("subscriber_id", "users(id)", "CASCADE", "CASCADE")
  con.Model(&models.NewArticlesSubscriber{}).AddForeignKey("article_id", "articles(id)", "CASCADE", "CASCADE")
}

func ImportDemoFixtures(c *Config)  {
  con := ConnectToDatabase(c)
  fixtures := getFixtures()
  tx := con.Begin()
  for _, user := range fixtures["Users"] {
    var newUser models.User
    newUser.ID = uint(user["id"].(int))
    newUser.Name = user["name"].(string)
    newUser.Email = user["email"].(string)
    newUser.Password = user["password"].(string)
    if err := tx.Create(&newUser).Error; err != nil {
      tx.Rollback()
      fmt.Println("Error while importing fixtures")
      return
    }
  }
  for _, article := range fixtures["Articles"] {
    var newArticle models.Article
    newArticle.Title = article["title"].(string)
    newArticle.Content = article["content"].(string)
    newArticle.UserID = uint(article["userId"].(int))
    if err := tx.Create(&newArticle).Error; err != nil {
      tx.Rollback()
      fmt.Println("Error while importing fixtures")
      return
    }
  }
  tx.Commit()
  fmt.Println("All fixtures have been imported successfully")
}

/**
 * Get fixture object
 */
func getFixtures() map[string][]map[string]interface{} {
  return map[string][]map[string]interface{} {
      "Users": []map[string]interface{} {
        {
          "id": 1,
          "name": "Stan",
          "email": "stan@gmail.com",
          "password": "40bd001563085fc35165329ea1ff5c5ecbdbbeef",
        },
        {
          "id": 2,
          "name": "Abrasha",
          "email": "abrasha@gmail.com",
          "password": "40bd001563085fc35165329ea1ff5c5ecbdbbeef",
        },
      },
      "Articles": []map[string]interface{} {
        {
          "title": "The 6 Strangest Presidential Elections in US History",
          "content": `Political news would have you thinking the 2016 presidential election was the nastiest, most contentious and most important our nation ever faced. However, in the annals of American elections, that one barely registers at least for sheer strangeness.

In fact, electoral politics have always been a down-and-dirty business, starting at least as early as 1800, when our founding fathers proved themselves adept at bitter battles. Other elections have featured nasty accusations, bizarre happenstance and even the death of one of the candidates.

Read on for six of the strangest presidential elections in U.S. history.`,
          "userId": 1,
        },
        {
          "title": "Why Toddlers Are So Bad at Hide-and-Seek",
          "content": `Young children across the globe enjoy playing games of hide and seek. There's something highly exciting for children about escaping someone else's glance and making oneself "invisible."

However, developmental psychologists and parents alike continue to witness that before school age, children are remarkably bad at hiding. Curiously, they often cover only their face or eyes with their hands, leaving the rest of their bodies visibly exposed.

For a long time, this ineffective hiding strategy was interpreted as evidence that young children are hopelessly "egocentric" creatures. Psychologists theorized that preschool children cannot distinguish their own perspective from someone else's. Conventional wisdom held that, unable to transcend their own viewpoint, children falsely assume that others see the world the same way they themselves do. So psychologists assumed children "hide" by covering their eyes because they conflate their own lack of vision with that of those around them.`,
          "userId": 1,
        },
        {
          "title": "Pope Extends Forgiveness for Abortion",
          "content": `The pope has extended the Catholic Church's forgiveness for abortion indefinitely.

In an apostolic letter dated Nov. 20, Pope Francis officially changed the church practice so that it now allows any parish priest to hear confessions from those who have obtained or performed an abortion, and to offer absolution.

The move comes as the Year of Mercy, or the Extraordinary Jubilee, comes to a close. With roots in the Old Testament, every 50 years, a jubilee year was designated as a time of forgiveness, as a reminder of God's mercy, according to Jewish tradition. The Catholic Church calls one every 25 years, and Pope Francis designated Dec. 8, 2015, to Nov. 20, 2016, as the Extraordinary Jubilee. During this time, the pope hoped followers would direct their attention on mercy so that we may become a more effective sign of the Father's action in our lives.`,
          "userId": 2,
        },
        {
          "title": "New 'Science Comics' Books Tackle Sharks, Brains, Drones and More",
          "content": `Once upon a time, comics were primarily the domain of costumed heroines and heroes preoccupied with battling evil supervillains and saving the planet. But generations of comics creators have proven that the graphic format used by comics can convey a wide range of narratives.

And a series of nonfiction graphic novels is proving that comics are terrific for telling stories about science.

Today (Nov. 16), First Second Books announced 13 upcoming titles in their "Science Comics" book series, to be released from 2017 through 2019. These nonfiction graphic novels combine engaging and vibrant artwork with characters who will introduce readers to a range of fascinating science topics: the history of drones, the evolution of the human brain, crows' intelligence, and the unexpectedly compelling life of trees.

The notion of combining graphic storytelling with science subjects made sense to First Second — in a statement, they called Science Comics an "amazing mix of two nerdy things that go excellently together." The series launched in March 2016 with two volumes: "Coral Reefs: Cities of the Ocean," and "Dinosaurs: Fossils and Feathers," followed by "Volcanoes: Fire and Life," published Nov. 15, and "Bats: Learning to Fly," which will be available Feb. 28, 2017.`,
          "userId": 1,
        },
        {
          "title": "Winning Data Visualizations Reveal Information Is Beautiful",
          "content": `From a global temperature timeline dating back to the last ice age, to an evolutionary history of changing tastes in pop music, some of the past year's most outstanding data visualizations were recognized on Nov. 2, at the Kantar "Information Is Beautiful" awards.

The annual awards — launched in 2012 by brand-research company Kantar and data journalist David McCandless, author of "Information is Beautiful" (Collins, 2000) — highlight exceptional examples of data visualization. This transformation of data into images — abstract or representational — relays information visually, to communicate a story within the numbers.

Winners in 2016 included visualizations of the amount of dialogue spoken by men and women in nearly four decades of movies, global predictions for seasonal winds, and the deadliest mass shootings in the U.S.`,
          "userId": 1,
        },
        {
          "title": "Why Humans Don't Have More Neanderthal Genes",
          "content": `Neanderthals and modern humans interbred long ago, but evolution has purged many of our caveman relative's genes from modern human genomes, a new study finds.

Neanderthals were the closest extinct relatives of modern humans. Previous research suggested that modern humans migrating out of Africa encountered and interbred with Neanderthals tens of thousands of years ago.

"We know that the ancestors of modern Europeans and Asians mated with Neanderthals, and as a result, the modern-day descendants of those people have some small amount of Neanderthal DNA in their genomes," said study lead author Ivan Juric, an evolutionary biologist at the University of California, Davis.`,
          "userId": 2,
        },
    },
  }
}
