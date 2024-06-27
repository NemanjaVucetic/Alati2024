import pytest
import time


from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.chrome.webdriver import WebDriver
from selenium.webdriver.common.action_chains import ActionChains
from selenium.webdriver.support.ui import Select
import chromedriver_autoinstaller


@pytest.fixture()
def browser():
    # instaliraj chrome driver
    # ukoliko ne postoji
    if chromedriver_autoinstaller.get_chrome_version() is None:
        chromedriver_autoinstaller.install()
    # neophodno odraditi ovo zbog problema sa CORS
    options = webdriver.ChromeOptions()
    options.add_argument("--disable-web-security")
    options.add_argument("--disable-site-isolation-trials")
    # kreiraj instancu drajvera za test
    driver = webdriver.Chrome(options=options)

    driver.implicitly_wait(10)
    # Yield the WebDriver instance
    yield driver
    # Close the WebDriver instance
    driver.quit()


class Test:

    def test_getAllGroup(self , browser:'WebDriver'):
        browser.get('http://localhost:8080/swagger/index.html#/')
        get_btn = browser.find_element(By.CSS_SELECTOR,"#operations-configGroups-get_configGroups_ > div > button > span.opblock-summary-method")
        get_btn.click()
        btn2 = browser.find_element(By.XPATH , "//*[@id='operations-configGroups-get_configGroups_']/div[2]/div/div[2]/div[1]/div[2]/button")
        btn2.click()
        btn3 = browser.find_element(By.XPATH , "//*[@id='operations-configGroups-get_configGroups_']/div[2]/div/div[3]/button")
        btn3.click()

        result = browser.find_element(By.XPATH , "//*[text()='200']")
        assert result.is_displayed()
        assert result.text == "200"
    
    def test_getAllConf(self , browser:'WebDriver'):
        browser.get('http://localhost:8080/swagger/index.html#/')
        get_btn = browser.find_element(By.CSS_SELECTOR,"#operations-configs-get_configs_ > div > button > span.opblock-summary-method")
        get_btn.click()
        btn2 = browser.find_element(By.XPATH , "//*[@id='operations-configs-get_configs_']/div[2]/div/div[2]/div[1]/div[2]/button")
        btn2.click()
        btn3 = browser.find_element(By.XPATH , "//*[@id='operations-configs-get_configs_']/div[2]/div/div[3]/button")
        btn3.click()

        result = browser.find_element(By.XPATH , "//*[text()='200']")
        assert result.is_displayed()
        assert result.text == "200"

    def test_getOneConf(self , browser:'WebDriver'):
        browser.get('http://localhost:8080/swagger/index.html#/')
        get_btn = browser.find_element(By.CSS_SELECTOR,"#operations-configs-get_configs__name___version_ > div > button > span.opblock-summary-method")
        get_btn.click()
        btn2 = browser.find_element(By.XPATH , "//*[@id='operations-configs-get_configs__name___version_']/div[2]/div/div[2]/div[1]/div[2]/button")
        btn2.click()
        input1 = browser.find_element(By.XPATH , "//*[@id='operations-configs-get_configs__name___version_']/div[2]/div/div[2]/div[2]/div/table/tbody/tr[1]/td[2]/input")
        input1.send_keys(Keys.BACK_SPACE * 100)
        input1.send_keys("db_config")
        input2 = browser.find_element(By.XPATH , "//*[@id='operations-configs-get_configs__name___version_']/div[2]/div/div[2]/div[2]/div/table/tbody/tr[2]/td[2]/input")
        input2.send_keys(Keys.BACK_SPACE * 100)
        input2.send_keys("2")
        btn3 = browser.find_element(By.XPATH , "//*[@id='operations-configs-get_configs__name___version_']/div[2]/div/div[3]/button[1]")
        btn3.click()

        result = browser.find_element(By.XPATH , "//*[text()='200']")
        assert result.is_displayed()
        assert result.text == "200"

    def test_getOneConfFail(self , browser:'WebDriver'):
        browser.get('http://localhost:8080/swagger/index.html#/')
        get_btn = browser.find_element(By.CSS_SELECTOR,"#operations-configs-get_configs__name___version_ > div > button > span.opblock-summary-method")
        get_btn.click()
        btn2 = browser.find_element(By.XPATH , "//*[@id='operations-configs-get_configs__name___version_']/div[2]/div/div[2]/div[1]/div[2]/button")
        btn2.click()
        input1 = browser.find_element(By.XPATH , "//*[@id='operations-configs-get_configs__name___version_']/div[2]/div/div[2]/div[2]/div/table/tbody/tr[1]/td[2]/input")
        input1.send_keys(Keys.BACK_SPACE * 100)
        input1.send_keys("db_config")
        input2 = browser.find_element(By.XPATH , "//*[@id='operations-configs-get_configs__name___version_']/div[2]/div/div[2]/div[2]/div/table/tbody/tr[2]/td[2]/input")
        input2.send_keys(Keys.BACK_SPACE * 100)
        input2.send_keys("1")
        btn3 = browser.find_element(By.XPATH , "//*[@id='operations-configs-get_configs__name___version_']/div[2]/div/div[3]/button[1]")
        btn3.click()

        result = browser.find_element(By.XPATH , "//*[text()='Failed to fetch.']")
        assert result.is_displayed()
        assert result.text == "Failed to fetch."

    def test_getOneGroup(self , browser:'WebDriver'):
        browser.get('http://localhost:8080/swagger/index.html#/')
        get_btn = browser.find_element(By.CSS_SELECTOR,"#operations-configGroups-get_configGroups__name___version_ > div.opblock-summary.opblock-summary-get > button > span.opblock-summary-method")
        get_btn.click()
        btn2 = browser.find_element(By.XPATH , "//*[@id='operations-configGroups-get_configGroups__name___version_']/div[2]/div/div[2]/div[1]/div[2]/button")
        btn2.click()
        input1 = browser.find_element(By.XPATH , "//*[@id='operations-configGroups-get_configGroups__name___version_']/div[2]/div/div[2]/div[2]/div/table/tbody/tr[1]/td[2]/input")
        input1.send_keys(Keys.BACK_SPACE * 100)
        input1.send_keys("db_cg")
        input2 = browser.find_element(By.XPATH , "//*[@id='operations-configGroups-get_configGroups__name___version_']/div[2]/div/div[2]/div[2]/div/table/tbody/tr[2]/td[2]/input")
        input2.send_keys(Keys.BACK_SPACE * 100)
        input2.send_keys("2")
        btn3 = browser.find_element(By.XPATH , "//*[@id='operations-configGroups-get_configGroups__name___version_']/div[2]/div/div[3]/button")
        btn3.click()

        result = browser.find_element(By.XPATH , "//*[text()='200']")
        assert result.is_displayed()
        assert result.text == "200"