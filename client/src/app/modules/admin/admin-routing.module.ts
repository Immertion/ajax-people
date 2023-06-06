import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { UnmoderatedPostsComponent } from "./unmoderated-posts/unmoderated-posts.component";
import { TagsPageComponent } from "./tags-page/tags-page.component";
import { AdminUsersComponent } from "./admin-users/admin-users.component";

const routes: Routes = [
  {
    path: 'tags',
    component: TagsPageComponent
  },
  {
    path: 'posts',
    component: UnmoderatedPostsComponent
  },
  {
    path: 'users',
    component: AdminUsersComponent
  }

];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class AdminRoutingModule { }
