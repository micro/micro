export class Page<T> {
  public pageSize = 10;
  public pageNo = 1;
  public items: T[] = [];
  public total: number = 0;
}
