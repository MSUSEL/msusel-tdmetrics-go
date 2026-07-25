namespace GroundTruth;

public readonly record struct Modifiers(int Value) {
    public bool None => this.Value == 0;
    private bool flag(int offset) => (this.Value & (1 << offset)) != 0;
    public bool Public       => this.flag(0);
    public bool Private      => this.flag(1);
    public bool Protected    => this.flag(2);
    public bool Static       => this.flag(3);
    public bool Final        => this.flag(4);
    public bool Synchronized => this.flag(5);
    public bool Volatile     => this.flag(6);
    public bool Transient    => this.flag(7);
    public bool Native       => this.flag(8);
    public bool Sealed       => this.flag(9);
    public bool Abstract     => this.flag(10);
    public bool StrictFP     => this.flag(11);
    public bool NonSealed    => this.flag(12);
    public bool Module       => this.flag(13);
    public bool Default      => this.flag(14);
}
