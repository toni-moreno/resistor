import { ResistorPage } from './app.po';

describe('Resistor App', function() {
  let page: ResistorPage;

  beforeEach(() => {
    page = new ResistorPage();
  });

  it('should display message saying app works', () => {
    page.navigateTo();
    expect(page.getParagraphText()).toEqual('app works!');
  });
});
