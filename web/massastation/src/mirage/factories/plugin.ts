import { Factory } from 'miragejs';
import { faker } from '@faker-js/faker';
import { IMassaPlugin } from '../../../../shared/interfaces/IPlugin';

export const pluginFactory = Factory.extend<IMassaPlugin>({
  id() {
    return faker.number.int().toString();
  },
  name() {
    return faker.lorem.word();
  },
  author() {
    // there is a 30% chance that the author will be MassaLabs
    return Math.random() < 0.3 ? 'Massa Labs' : faker.person.firstName();
  },
  description() {
    return faker.lorem.sentence();
  },
  home() {
    const name: string = this.name.toString().toLowerCase();
    return `/plugin/massalabs/${name}/`;
  },
  logo() {
    const name = this.name.toString().toLowerCase();
    return `/plugin/massalabs/${name}/logo.svg`;
  },
  status: 'Up',
  updatable() {
    return Math.random() < 0.5;
  },
});
