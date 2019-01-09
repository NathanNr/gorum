import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Title } from '@angular/platform-browser';
import { Config } from '../config';
import { Thread } from '../thread/thread.component';

export class Board {
  id: number;
  name: string;
  description: string;
  icon: string;
  sort: number;

  constructor(id: number, name: string, description: string, icon: string, sort: number) {
    this.id = id;
    this.name = name;
    this.description = description;
    this.icon = icon;
    this.sort = sort;
  }
}

@Component({
  selector: 'app-board',
  templateUrl: './board.component.html',
  styleUrls: ['./board.component.css']
})
export class BoardComponent implements OnInit {
  config = Config;
  conf = Config.get;
  lang = Config.lang;

  threads: Thread[] = [];
  id = +this.route.snapshot.paramMap.get('id');

  constructor(private route: ActivatedRoute, private title: Title) { }

  ngOnInit() {
    Config.setLogin(false);
    Config.API('board', { boardID: this.id }).subscribe(values => this.title.setTitle(values['name'] + ' - ' + Config.get('title')));
    Config.API('threads', { boardID: this.id }).subscribe(values => this.listThreads(values));
  }

  listThreads(values: any) {
    Object.entries(values).forEach(thread => this.threads.push(new Thread(<number>thread[1]['id'], <string>thread[1]['name'],
      <string>thread[1]['board'], <number>thread[1]['author'], <number>thread[1]['created'], <string>thread[1]['content'],
      <string>thread[1]['authorName'], <string>thread[1]['authorAvatar'], null)));
    this.threads.sort((a, b) => b.created - a.created);
  }
}
