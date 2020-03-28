import scrapy
from scrapy.http import FormRequest

from ..items import QuoteItem


class QuoteSpider(scrapy.Spider):
    name = "Quotes"
    start_urls = ['http://quotes.toscrape.com/login']
    def parse(self, response):
        csrf = response.css('form input::attr(value)').extract_first()
        return FormRequest.from_response(response, formdata={
            'csrf_token':csrf,
            'username':'xyz',
            'password':'xxx'
        }, callback=self.start_scrapping)

    def start_scrapping(self, response):
        title = response.css('title::text')[0].extract()
        all_div_quotes = response.css("div.quote")
        for quote_box in all_div_quotes:
            quote = quote_box.css("span.text::text").extract_first()
            author = quote_box.css(".author::text").extract_first()
            tags = quote_box.css(".tag::text").extract()
            item = QuoteItem()
            item['author'] = author
            item['quote'] = quote
            item['tags'] = tags
            yield item
        print("page done")
        next_page = response.css('li.next a::attr(href)').get()
        if next_page:
            yield response.follow(next_page, self.start_scrapping)
